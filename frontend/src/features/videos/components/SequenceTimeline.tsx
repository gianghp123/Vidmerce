import { SequenceItemComponent } from "@/components/custom/SequenceItem";
import { PlusCircle, Wand2 } from "lucide-react";
import type { SequenceItem } from "../models/video-sequence.model";

interface SequenceTimelineProps {
  sequences: SequenceItem[];
  onSequencesChange: (sequences: SequenceItem[]) => void;
  availableAssets?: { id: string; name: string; imageUrl: string }[];
  className?: string;
}

export function SequenceTimeline({
  sequences,
  onSequencesChange,
  availableAssets = [],
}: SequenceTimelineProps) {
  const handleAddSequence = () => {
    const newSequence: SequenceItem = {
      id: Math.random().toString(36).substring(7),
      assetId: "",
      duration: 3,
      transition: "FADE",
    };
    onSequencesChange([...sequences, newSequence]);
  };

  const handleUpdateSequence = (updatedItem: SequenceItem) => {
    const newSequences = sequences.map((seq) =>
      seq.id === updatedItem.id ? updatedItem : seq
    );
    onSequencesChange(newSequences);
  };

  const handleRemoveSequence = (id: string) => {
    onSequencesChange(sequences.filter((seq) => seq.id !== id));
  };

  const handleAutoArrange = () => {
    if (sequences.length <= 1) return;
    const newSequences = [...sequences].sort(() => Math.random() - 0.5);
    onSequencesChange(newSequences);
  };

  return (
    <section>
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-2xl font-bold text-on-surface">Sequence Timeline</h3>
        <button
          onClick={handleAutoArrange}
          className="text-sm font-bold text-secondary flex items-center gap-1.5 hover:opacity-80 transition-opacity"
        >
          <Wand2 className="w-4 h-4" />
          Auto-Arrange
        </button>
      </div>

      <div className="space-y-4">
        {/* Assumes SequenceItemComponent uses white bg and rounded-2xl */}
        {sequences.map((sequence) => (
          <SequenceItemComponent key={sequence.id} item={sequence} onChange={handleUpdateSequence} onRemove={handleRemoveSequence} assets={availableAssets} />
        ))}

        <button
          onClick={handleAddSequence}
          className="w-full border-2 border-dashed border-outline-variant/40 h-32 rounded-[24px] flex flex-col items-center justify-center text-on-surface-variant hover:bg-surface-container-low/40 transition-all"
        >
          <PlusCircle className="w-8 h-8 mb-2 text-on-surface-variant/40" />
          <span className="text-[15px] font-bold text-on-surface-variant/80">Add Sequence Item</span>
        </button>
      </div>
    </section>
  );
}
