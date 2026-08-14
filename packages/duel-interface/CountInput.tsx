interface CountInputProps {
  value: number;
  onChange: (value: number) => void;
  min?: number;
  max?: number;
}

export function CountInput({
  value,
  onChange,
  min = 0,
  max = Infinity,
}: CountInputProps) {
  const handleDecrement = () => {
    if (value > min) {
      onChange(value - 1);
    }
  };

  const handleIncrement = () => {
    if (value < max) {
      onChange(value + 1);
    }
  };

  const isAtMin = value <= min;
  const isAtMax = value >= max;

  return (
    <div className="inline-flex items-center gap-1">
      <button
        onClick={handleDecrement}
        disabled={isAtMin}
        className={`
          w-8 h-8 rounded-lg border-2 flex items-center justify-center
          font-semibold text-lg transition-colors
          ${
            isAtMin
              ? "border-gray-700 bg-gray-800 text-gray-600 cursor-not-allowed"
              : "border-gray-600 bg-gray-700 text-white hover:bg-gray-600 hover:border-gray-500 active:bg-gray-500"
          }
        `}
        aria-label="Decrease"
      >
        −
      </button>

      <div className="w-8 h-8 rounded-lg border-2 border-gray-600 bg-gray-800 flex items-center justify-center font-semibold text-sm text-white">
        {value}
      </div>

      <button
        onClick={handleIncrement}
        disabled={isAtMax}
        className={`
          w-8 h-8 rounded-lg border-2 flex items-center justify-center
          font-semibold text-lg transition-colors
          ${
            isAtMax
              ? "border-gray-700 bg-gray-800 text-gray-600 cursor-not-allowed"
              : "border-gray-600 bg-gray-700 text-white hover:bg-gray-600 hover:border-gray-500 active:bg-gray-500"
          }
        `}
        aria-label="Increase"
      >
        +
      </button>
    </div>
  );
}
