import { useEffect, useState } from "react";
import { Button } from "../common/Button";
import { FieldInput } from "./FieldInput";

interface SaveDraftDialogProps {
  error?: string | null;
  onCancel: () => void;
  onConfirm: (name: string, note: string) => Promise<void> | void;
  open: boolean;
  saving: boolean;
}

export function SaveDraftDialog({
  error = null,
  onCancel,
  onConfirm,
  open,
  saving,
}: SaveDraftDialogProps) {
  const [name, setName] = useState("");
  const [note, setNote] = useState("");

  useEffect(() => {
    if (!open) {
      setName("");
      setNote("");
    }
  }, [open]);

  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/30 px-4">
      <div
        aria-labelledby="save-draft-title"
        aria-modal="true"
        className="w-full max-w-lg rounded-3xl bg-white p-6 shadow-2xl"
        role="dialog"
      >
        <div className="space-y-2">
          <h2
            className="text-xl font-semibold text-slate-900"
            id="save-draft-title"
          >
            保存为草稿
          </h2>
          <p className="text-sm text-slate-500">
            主简历不会被覆盖，当前编辑内容会单独保存为一个版本。
          </p>
        </div>
        <div className="mt-6 space-y-4">
          <FieldInput
            label="草稿名称"
            onChange={setName}
            placeholder="例如：面试前调整版"
            value={name}
          />
          <FieldInput
            label="草稿备注"
            onChange={setNote}
            placeholder="补充这次保存的变更说明"
            textarea
            value={note}
          />
          {error ? <p className="text-sm text-rose-600">{error}</p> : null}
        </div>
        <div className="mt-6 flex items-center justify-end gap-3">
          <Button onClick={onCancel} variant="secondary">
            取消
          </Button>
          <Button
            onClick={() => onConfirm(name.trim(), note.trim())}
            variant="primary"
          >
            {saving ? "保存中..." : "确认保存草稿"}
          </Button>
        </div>
      </div>
    </div>
  );
}
