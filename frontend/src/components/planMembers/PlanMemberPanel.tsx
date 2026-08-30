"use client";

import { useState } from "react";
import DeleteIcon from "@mui/icons-material/Delete";
import { Alert, Button, CircularProgress, IconButton, TextField } from "@mui/material";

import { usePlanMemberContainer } from "@/containers/PlanMemberContainer";

type Props = { planId: string };

export function PlanMemberPanel({ planId }: Props) {
  const { memberVMs, loading, error, onAddMember, onRemoveMember } = usePlanMemberContainer(planId);
  const [userId, setUserId] = useState("");

  const handleAdd = async () => {
    const trimmed = userId.trim();
    if (!trimmed) return;
    const added = await onAddMember(trimmed);
    if (added) setUserId("");
  };

  return (
    <div className="flex w-full max-w-sm flex-col gap-4">
      <h2 className="text-lg font-semibold">メンバー</h2>

      <div className="flex items-start gap-2">
        <TextField
          size="small"
          label="ユーザーID"
          helperText="メールアドレスでの招待は未対応です"
          value={userId}
          onChange={(event) => setUserId(event.target.value)}
          disabled={loading}
          fullWidth
        />
        <Button variant="contained" onClick={handleAdd} disabled={loading || !userId.trim()}>
          追加
        </Button>
      </div>

      {loading && <CircularProgress aria-label="読み込み中" size={24} />}

      {!loading && error && <Alert severity="error">{error}</Alert>}

      <ul className="flex flex-col gap-2">
        {memberVMs.map((member) => (
          <li
            key={member.id}
            className="flex items-center justify-between rounded border p-2 text-sm"
          >
            <span>{member.userId}</span>
            <IconButton
              aria-label="削除"
              size="small"
              onClick={() => onRemoveMember(member.userId)}
              disabled={loading}
            >
              <DeleteIcon fontSize="small" />
            </IconButton>
          </li>
        ))}
      </ul>
    </div>
  );
}
