"use client";

import { useCallback, useEffect, useState } from "react";

import type { CreatePinInput, UpdatePinInput } from "@/domain/pins/types";
import { useNotes } from "@/hooks/useNotes";
import { usePlan } from "@/hooks/usePlan";
import { usePins } from "@/hooks/usePins";
import {
  toNoteViewModel,
  toPinViewModel,
  toPlanViewModel,
  type NoteViewModel,
  type PinViewModel,
  type PlanViewModel,
} from "@/mappers/planMapper";

type UsePlanContainerReturn = {
  planVM: PlanViewModel | null;
  loading: boolean;
  error: string | null;
  onCreatePin: (data: CreatePinInput) => Promise<void>;
  onUpdatePin: (pinId: string, data: UpdatePinInput) => Promise<void>;
  onDeletePin: (pinId: string) => Promise<void>;
  applyPinUpdate: (pin: PinViewModel) => void;
  onCreateNote: (pinId: string, content: string) => Promise<void>;
  onUpdateNote: (pinId: string, noteId: string, content: string) => Promise<void>;
  onDeleteNote: (pinId: string, noteId: string) => Promise<void>;
};

export function usePlanContainer(shareToken: string): UsePlanContainerReturn {
  const { getPlan, loading: planLoading, error: planError } = usePlan();
  const { createPin, updatePin, deletePin, loading: pinLoading, error: pinError } = usePins();
  const { createNote, updateNote, deleteNote, loading: noteLoading, error: noteError } = useNotes();

  const [planVM, setPlanVM] = useState<PlanViewModel | null>(null);

  useEffect(() => {
    getPlan(shareToken).then((plan) => {
      if (plan) setPlanVM(toPlanViewModel(plan));
    });
  }, [shareToken, getPlan]);

  const onCreatePin = useCallback(
    async (data: CreatePinInput) => {
      if (!planVM) return;
      const pin = await createPin(planVM.id, data);
      if (!pin) return;
      setPlanVM((prev) => (prev ? { ...prev, pins: [...prev.pins, toPinViewModel(pin)] } : prev));
    },
    [planVM, createPin],
  );

  const onUpdatePin = useCallback(
    async (pinId: string, data: UpdatePinInput) => {
      if (!planVM) return;
      const pin = await updatePin(planVM.id, pinId, data);
      if (!pin) return;
      setPlanVM((prev) =>
        prev
          ? {
              ...prev,
              pins: prev.pins.map((p) =>
                p.id === pinId ? { ...toPinViewModel(pin), notes: p.notes } : p,
              ),
            }
          : prev,
      );
    },
    [planVM, updatePin],
  );

  const onDeletePin = useCallback(
    async (pinId: string) => {
      if (!planVM) return;
      const ok = await deletePin(planVM.id, pinId);
      if (!ok) return;
      setPlanVM((prev) =>
        prev ? { ...prev, pins: prev.pins.filter((p) => p.id !== pinId) } : prev,
      );
    },
    [planVM, deletePin],
  );

  const applyPinUpdate = useCallback((pin: PinViewModel) => {
    setPlanVM((prev) =>
      prev
        ? {
            ...prev,
            pins: prev.pins.map((p) => (p.id === pin.id ? { ...pin, notes: p.notes } : p)),
          }
        : prev,
    );
  }, []);

  const onCreateNote = useCallback(
    async (pinId: string, content: string) => {
      if (!planVM) return;
      const note = await createNote(planVM.id, pinId, content);
      if (!note) return;
      setPlanVM((prev) =>
        prev
          ? {
              ...prev,
              pins: prev.pins.map((p) =>
                p.id === pinId ? { ...p, notes: [...p.notes, toNoteViewModel(note)] } : p,
              ),
            }
          : prev,
      );
    },
    [planVM, createNote],
  );

  const onUpdateNote = useCallback(
    async (pinId: string, noteId: string, content: string) => {
      if (!planVM) return;
      const note = await updateNote(planVM.id, pinId, noteId, content);
      if (!note) return;
      setPlanVM((prev) =>
        prev
          ? {
              ...prev,
              pins: prev.pins.map((p) =>
                p.id === pinId
                  ? {
                      ...p,
                      notes: p.notes.map((n) => (n.id === noteId ? toNoteViewModel(note) : n)),
                    }
                  : p,
              ),
            }
          : prev,
      );
    },
    [planVM, updateNote],
  );

  const onDeleteNote = useCallback(
    async (pinId: string, noteId: string) => {
      if (!planVM) return;
      const ok = await deleteNote(planVM.id, pinId, noteId);
      if (!ok) return;
      setPlanVM((prev) =>
        prev
          ? {
              ...prev,
              pins: prev.pins.map((p) =>
                p.id === pinId ? { ...p, notes: p.notes.filter((n) => n.id !== noteId) } : p,
              ),
            }
          : prev,
      );
    },
    [planVM, deleteNote],
  );

  return {
    planVM,
    loading: planLoading || pinLoading || noteLoading,
    error: planError ?? pinError ?? noteError,
    onCreatePin,
    onUpdatePin,
    onDeletePin,
    applyPinUpdate,
    onCreateNote,
    onUpdateNote,
    onDeleteNote,
  };
}

export type { PinViewModel, NoteViewModel, PlanViewModel };
