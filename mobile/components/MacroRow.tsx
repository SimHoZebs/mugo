import { MacroDisplay } from "@/components/MacroDisplay";

interface MacroRowProps {
  label: string;
  value: number;
  unit: string;
  colorClass: string;
}

export function MacroRow(props: MacroRowProps) {
  return <MacroDisplay variant="row" {...props} />;
}
