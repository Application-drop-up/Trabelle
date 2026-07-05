import { http, HttpResponse } from "msw";

import type { Plan } from "@/domain/plans/types";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const mockPlan: Plan = {
  id: "mock-plan-id",
  share_token: "mock-token",
  title: "東京旅行 2025",
  pins: [
    {
      id: "pin-1",
      plan_id: "mock-plan-id",
      name: "東京タワー",
      latitude: 35.6586,
      longitude: 139.7454,
      category: "sightseeing",
      colour: "#3B82F6",
      notes: [
        {
          id: "note-1",
          pin_id: "pin-1",
          content: "展望台から富士山が見えるかも",
          created_at: "2025-01-01T09:00:00Z",
          updated_at: "2025-01-01T09:00:00Z",
        },
      ],
      created_at: "2025-01-01T09:00:00Z",
      updated_at: "2025-01-01T09:00:00Z",
    },
    {
      id: "pin-2",
      plan_id: "mock-plan-id",
      name: "銀座 久兵衛",
      latitude: 35.6717,
      longitude: 139.7652,
      category: "restaurant",
      colour: "#EF4444",
      notes: [],
      created_at: "2025-01-01T12:00:00Z",
      updated_at: "2025-01-01T12:00:00Z",
    },
    {
      id: "pin-3",
      plan_id: "mock-plan-id",
      name: "渋谷スクランブル交差点",
      latitude: 35.6595,
      longitude: 139.7004,
      category: "sightseeing",
      colour: "#3B82F6",
      notes: [
        {
          id: "note-2",
          pin_id: "pin-3",
          content: "夜景が綺麗",
          created_at: "2025-01-01T18:00:00Z",
          updated_at: "2025-01-01T18:00:00Z",
        },
      ],
      created_at: "2025-01-01T18:00:00Z",
      updated_at: "2025-01-01T18:00:00Z",
    },
  ],
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-01T00:00:00Z",
};

export const handlers = [
  http.get(`${API_BASE}/plans/:shareToken`, () => {
    return HttpResponse.json(mockPlan);
  }),

  http.post(`${API_BASE}/plans`, () => {
    return HttpResponse.json(mockPlan, { status: 201 });
  }),
];
