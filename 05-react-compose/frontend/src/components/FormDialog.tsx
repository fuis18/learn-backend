import { createContext, useContext, useState } from "react";
import type { ReactElement, ReactNode } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import AlertDialog from "@/components/AlertDialog";
import { useAlert, type AlertState } from "@/hooks/useAlert";

interface FormDialogContextValue {
  close: () => void;
  show: (alert: Omit<AlertState, "open">) => void;
}

const FormDialogContext = createContext<FormDialogContextValue | null>(null);

export function useFormDialog() {
  const ctx = useContext(FormDialogContext);
  if (!ctx) {
    throw new Error("useFormDialog must be used within a <FormDialog>");
  }
  return ctx;
}

interface FormDialogProps {
  trigger: ReactElement;
  title: ReactNode;
  children: ReactNode;
  className?: string;
  alertActionLabel?: string;
}

const FormDialog = ({
  trigger,
  title,
  children,
  className = "max-w-md",
  alertActionLabel = "Close",
}: FormDialogProps) => {
  const [open, setOpen] = useState(false);
  const { alert, show, hide } = useAlert();

  const value = {
    close: () => setOpen(false),
    show,
  };

  return (
    <FormDialogContext.Provider value={value}>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger render={trigger} />

        <DialogContent className={className}>
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
          </DialogHeader>

          {children}
        </DialogContent>

        <AlertDialog
          open={alert.open}
          onOpenChange={(next) => {
            if (!next) hide();
          }}
          variant={alert.variant}
          title={alert.title}
          description={alert.description}
          actionLabel={alertActionLabel}
        />
      </Dialog>
    </FormDialogContext.Provider>
  );
};

export default FormDialog;
