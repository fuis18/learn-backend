import { CircleAlert, CircleCheck } from "lucide-react";
import { cn } from "cn";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

export type AlertVariant = "success" | "error";

interface AlertDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  variant?: AlertVariant;
  actionLabel?: string;
}

const AlertDialog = ({
  open,
  onOpenChange,
  title,
  description,
  variant = "success",
  actionLabel = "Aceptar",
}: AlertDialogProps) => {
  const isSuccess = variant === "success";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader className="gap-2">
          <div className="flex items-center gap-2">
            <span
              className={cn(
                "flex size-8 shrink-0 items-center justify-center rounded-full",
                isSuccess
                  ? "bg-success/15 text-success"
                  : "bg-destructive/15 text-destructive",
              )}
            >
              {isSuccess ? (
                <CircleCheck className="size-4" />
              ) : (
                <CircleAlert className="size-4" />
              )}
            </span>
            <DialogTitle>{title}</DialogTitle>
          </div>
          {description && (
            <DialogDescription>{description}</DialogDescription>
          )}
        </DialogHeader>

        <DialogFooter>
          <DialogClose
            render={
              <Button
                type="button"
                variant={isSuccess ? "default" : "destructive"}
              />
            }
          >
            {actionLabel}
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default AlertDialog;
