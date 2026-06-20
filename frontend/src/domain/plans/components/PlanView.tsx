import { PinListItem } from "@/domain/pins/components/PinListItem";
import type { Plan } from "@/domain/plans/types";

export function PlanView({ plan }: { plan: Plan }) {
  return (
    <div className="flex flex-col gap-6 p-8">
      <h1 className="text-2xl font-semibold">{plan.title}</h1>
      <ul className="flex flex-col gap-4">
        {plan.pins.map((pin) => (
          <PinListItem key={pin.id} pin={pin} />
        ))}
      </ul>
    </div>
  );
}
