'use client';

import { ApiError, BaseResponse } from '@/types/base.model';
import { useActionState, useEffect, useRef, startTransition } from 'react';
import { toast } from 'sonner';

export interface UseServerActionOptions<T> {
  onSuccess?: (data: T, meta?: BaseResponse<T>['meta']) => void;
  onError?: (error: ApiError) => void;
  successMessage?: string;
  errorMessage?: string;
}

export function useServerAction<T, Payload = FormData>(
  action: (prevState: BaseResponse<T>, payload: Payload) => Promise<BaseResponse<T>>,
  options: UseServerActionOptions<T> = {},
  initialState: BaseResponse<T> = { data: null, error: null }
) {
  const [state, formAction, isPending] = useActionState(action, initialState);

  /**
   * ✅ Safe dispatcher for manual usage
   * Always wrapped in startTransition
   */
  const dispatch = (payload: Payload) => {
    startTransition(() => {
      formAction(payload);
    });
  };

  // Track previous pending state to detect transition (true -> false)
  const wasPending = useRef(false);

  // Avoid stale closures
  const optionsRef = useRef(options);
  useEffect(() => {
    optionsRef.current = options;
  }, [options]);

  useEffect(() => {
    // Only run when request just finished
    if (!wasPending.current || isPending) {
      wasPending.current = isPending;
      return;
    }

    const currentOptions = optionsRef.current;

    if (state.error) {
      // --- ERROR ---
      const errorMessage =
        currentOptions.errorMessage ||
        state.error.message ||
        'An error occurred';

      toast.error(errorMessage);
      currentOptions.onError?.(state.error);
    } else if (state.data != null) {
      // --- SUCCESS ---
      if (currentOptions.successMessage) {
        toast.success(currentOptions.successMessage);
      }

      currentOptions.onSuccess?.(state.data, state.meta);
    }

    wasPending.current = isPending;
  }, [isPending, state.error, state.data, state.meta]);

  return {
    state,
    formAction, // ✅ for <form action={formAction}>
    dispatch,   // ✅ for manual calls (safe)
    isPending,
  };
}