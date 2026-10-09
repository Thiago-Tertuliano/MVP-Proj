import * as React from "react";
import { Sparkles } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "@/components/ui/dialog";

export type LevelUpDialogProps = {
  open: boolean;
  nivel: number;
  onClose: () => void;
};

/**
 * Celebração de subida de nível. A animação só roda com `motion-safe`: quem pediu menos movimento
 * recebe o mesmo conteúdo, estático.
 */
export function LevelUpDialog({ open, nivel, onClose }: LevelUpDialogProps) {
  return (
    <Dialog open={open} onOpenChange={(aberto) => !aberto && onClose()}>
      <DialogContent className="max-w-sm items-center text-center">
        <div className="mx-auto flex h-24 w-24 items-center justify-center rounded-full bg-notify text-primary-foreground motion-safe:animate-in motion-safe:zoom-in-50 motion-safe:duration-500">
          <span className="text-4xl font-bold tabular-nums" aria-hidden="true">
            {nivel}
          </span>
        </div>
        <div className="space-y-1">
          <DialogTitle className="flex items-center justify-center gap-2 text-xl">
            <Sparkles className="h-5 w-5 text-notify" aria-hidden="true" />
            Subiu de nível!
          </DialogTitle>
          <DialogDescription>
            Você chegou ao <strong>nível {nivel}</strong>. Continue dominando nós para ir mais longe.
          </DialogDescription>
        </div>
        <Button type="button" onClick={onClose} className="mx-auto">
          Continuar
        </Button>
      </DialogContent>
    </Dialog>
  );
}
