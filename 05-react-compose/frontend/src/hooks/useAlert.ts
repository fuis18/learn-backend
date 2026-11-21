import { useState } from "react";
import type { AlertVariant } from "@/types/alert";

export interface AlertState {
  open: boolean;
  variant: AlertVariant;
  title: string;
  description?: string;
}

const IDLE_ALERT: AlertState = {
  open: false,
  variant: "success",
  title: "",
};

export function useAlert() {
  const [alert, setAlert] = useState<AlertState>(IDLE_ALERT);

  const show = (next: Omit<AlertState, "open">) =>
    setAlert({ ...next, open: true });

  const hide = () => setAlert((prev) => ({ ...prev, open: false }));

  return { alert, show, hide };
}
