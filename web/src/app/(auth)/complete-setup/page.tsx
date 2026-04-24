"use client";

import { Spinner } from "@/components/ui/spinner";
import { ROUTES } from "@/lib/routes";
import { useUser } from "@clerk/react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import Link from "next/link";
import { useRouter } from "next/router";

export default function CompleteSetupPage() {
  const { user } = useUser();
  const navigate = useRouter();

  const [retries, setRetries] = useState(0);
  const [showErrorModal, setShowErrorModal] = useState(false);

  const maxRetries = 15;

  useEffect(() => {
    if (!user) return;

    // stop retry when exceeded
    if (retries >= maxRetries) {
      setShowErrorModal(true);
      return;
    }

    const timeout = setTimeout(async () => {
      try {
        await user.reload();

        if (user.publicMetadata?.role) {
          navigate.push(ROUTES.MAIN.ASSETS.LIST);
          return;
        }
      } catch (error) {
        console.error("Role sync error:", error);
      }

      setRetries((prev) => prev + 1);
    }, 1000);

    return () => clearTimeout(timeout);
  }, [user, retries, navigate]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-8">
      <div className="text-center">
        <Spinner className="m-auto mb-4 h-12 w-12" />
        <p className="text-muted-foreground">
          Setting up your account...
        </p>
      </div>

      {/* Error Modal */}
      <Dialog open={showErrorModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Setup failed</DialogTitle>
            <DialogDescription>
              There is an error with your account setup. Please try again.
            </DialogDescription>
          </DialogHeader>

          <DialogFooter>
            <Button variant="outline" asChild>
              <Link href={ROUTES.MAIN.AUTH.SIGN_UP}>Sign Up Again</Link>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}