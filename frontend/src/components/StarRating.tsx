interface StarRatingProps {
  value: number;
  onChange?: (value: number) => void;
  readonly?: boolean;
  size?: "sm" | "md" | "lg";
}

export function StarRating({ value, onChange, readonly = false, size = "md" }: StarRatingProps) {
  const sizePx = { sm: "1.1rem", md: "1.6rem", lg: "2rem" }[size];

  return (
    <span style={{ display: "inline-flex", gap: "2px" }}>
      {[1, 2, 3, 4, 5].map((star) => (
        <span
          key={star}
          onClick={() => !readonly && onChange?.(star)}
          style={{
            fontSize: sizePx,
            cursor: readonly ? "default" : "pointer",
            color: star <= value ? "#c47a1b" : "#d4c5a9",
            userSelect: "none",
            transition: "color 0.15s",
          }}
        >
          ★
        </span>
      ))}
    </span>
  );
}
