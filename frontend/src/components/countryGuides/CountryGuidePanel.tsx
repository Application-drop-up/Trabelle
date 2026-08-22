"use client";

import { Alert, Chip, CircularProgress, MenuItem, TextField } from "@mui/material";

import { useCountryGuideContainer } from "@/containers/CountryGuideContainer";
import { COUNTRY_GUIDE_ITEM_CATEGORY_LABELS } from "@/domain/countryGuides/labels";

export function CountryGuidePanel() {
  const { guideVMs, selectedCode, selectedGuideVM, loading, error, onSelectCode } =
    useCountryGuideContainer();

  return (
    <div className="flex w-full max-w-sm flex-col gap-4">
      <h2 className="text-lg font-semibold">渡航ガイド</h2>

      <TextField
        select
        label="渡航先"
        value={selectedCode}
        onChange={(event) => onSelectCode(event.target.value)}
        disabled={loading && guideVMs.length === 0}
        fullWidth
      >
        {guideVMs.map((guideVM) => (
          <MenuItem key={guideVM.id} value={guideVM.countryCode}>
            {guideVM.countryName}
          </MenuItem>
        ))}
      </TextField>

      {loading && <CircularProgress aria-label="読み込み中" size={24} />}

      {!loading && error && <Alert severity="error">{error}</Alert>}

      {!loading && !error && selectedGuideVM && (
        <ul className="flex flex-col gap-3">
          {selectedGuideVM.items.map((item) => (
            <li key={item.id} className="flex flex-col gap-1 rounded border p-3 text-sm">
              <div className="flex items-center gap-2">
                <Chip label={COUNTRY_GUIDE_ITEM_CATEGORY_LABELS[item.category]} size="small" />
                {item.isMandatory && <Chip label="必須" size="small" color="error" />}
              </div>
              <span className="font-medium">{item.title}</span>
              {item.description && <span className="text-gray-600">{item.description}</span>}
              {item.url && (
                <a
                  href={item.url}
                  target="_blank"
                  rel="noreferrer"
                  className="text-blue-600 underline"
                >
                  詳細を見る
                </a>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
