export function parsePrevKeys(params: URLSearchParams): string[] {
  const prevRaw = params.get("prevKeys");
  if (!prevRaw) return [];
  try {
    const parsed = JSON.parse(prevRaw);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

export function buildForwardParams(
  currentParams: URLSearchParams,
  nextLastKey: string,
  limit: number
): URLSearchParams {
  const next = new URLSearchParams(currentParams.toString());
  const currentKey = currentParams.get("lastKey");
  const prevKeys = parsePrevKeys(currentParams);

  // Push CURRENT key to history before moving forward
  const newStack = currentKey ? [...prevKeys, currentKey] : prevKeys;

  next.set("lastKey", nextLastKey);
  next.set("limit", limit.toString());

  if (newStack.length > 0) {
    next.set("prevKeys", JSON.stringify(newStack));
  }
  return next;
}

export function buildBackParams(currentParams: URLSearchParams): URLSearchParams | null {
  const currentKey = currentParams.get("lastKey");
  const prevKeys = parsePrevKeys(currentParams);
  
  if (!currentKey) return null;

  const next = new URLSearchParams(currentParams.toString());

  if (prevKeys.length === 0) {
    // Going back to Page 1
    next.delete("lastKey");
    next.delete("prevKeys");
  } else {
    // Pop the last key from the stack
    const newStack = [...prevKeys];
    const targetKey = newStack.pop();

    if (targetKey) {
      next.set("lastKey", targetKey);
      if (newStack.length > 0) {
        next.set("prevKeys", JSON.stringify(newStack));
      } else {
        next.delete("prevKeys");
      }
    }
  }

  return next;
}