import { TemplateCard, TEMPLATES } from "./TemplateCard";

interface TemplateGridProps {
  selectedIndex: number;
  onSelectTemplate: (index: number) => void;
}

export function TemplateGrid({ selectedIndex, onSelectTemplate }: TemplateGridProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {TEMPLATES.map((template, index) => (
        <TemplateCard
          key={template.id}
          title={template.title}
          description={template.description}
          type={template.type}
          isSelected={selectedIndex === index}
          onClick={() => onSelectTemplate(index)}
        />
      ))}
    </div>
  );
}