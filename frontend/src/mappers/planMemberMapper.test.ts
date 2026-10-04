import type { PlanMember } from "@/domain/planMembers/types";
import { toPlanMemberViewModel } from "./planMemberMapper";

const mockMember: PlanMember = {
  id: "member-1",
  plan_id: "plan-1",
  user_id: "user-1",
  created_at: "2024-01-01T00:00:00Z",
};

describe("toPlanMemberViewModel", () => {
  it("converts snake_case PlanMember to camelCase ViewModel", () => {
    const vm = toPlanMemberViewModel(mockMember);

    expect(vm.id).toBe("member-1");
    expect(vm.planId).toBe("plan-1");
    expect(vm.userId).toBe("user-1");
    expect(vm.createdAt).toBe("2024-01-01T00:00:00Z");
  });
});
