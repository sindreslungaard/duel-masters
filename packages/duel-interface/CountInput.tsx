interface CountInputProps {
  value: number;
  onChange: (value: number) => void;
  min?: number;
  max?: number;
}

/** Roughly how many presses of a coarse button should cross the whole range.
 * Low enough that a wide range is quick to cross, high enough that the coarse
 * step stays a recognisable round number. */
const COARSE_PRESSES = 12;

/** Round numbers a coarse step is allowed to be. Powers land on 500, which is
 * the granularity every creature in the game is printed at. */
const COARSE_STEPS = [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000];

/** The jump a coarse button makes, worked out from the range alone: a prompt
 * for a power (0-6000) steps in 500s, while one for a number of cards is small
 * enough that ±1 is already the right size and no coarse button is offered. */
export function coarseStepFor(min: number, max: number): number {
  const range = max - min;

  if (!Number.isFinite(range) || range <= COARSE_PRESSES) {
    return 1;
  }

  const ideal = range / COARSE_PRESSES;

  return COARSE_STEPS.find((step) => step >= ideal) ?? range;
}

export function CountInput({
  value,
  onChange,
  min = 0,
  max = Infinity,
}: CountInputProps) {
  const coarseStep = coarseStepFor(min, max);
  const hasSlider = Number.isFinite(max) && max > min;

  const clamp = (n: number) => Math.min(max, Math.max(min, n));

  /** Coarse buttons land on round numbers rather than carrying an offset the
   * ±1 buttons left behind, so +500 from 2501 gives 3000, not 3001. */
  const stepBy = (delta: number) => {
    const moved = value + delta;

    if (Math.abs(delta) === 1) {
      onChange(clamp(moved));
      return;
    }

    const grid = min + Math.round((moved - min) / coarseStep) * coarseStep;

    onChange(clamp(grid));
  };

  const button = (delta: number, label: string) => {
    const disabled = delta < 0 ? value <= min : value >= max;

    return (
      <button
        type="button"
        onClick={() => stepBy(delta)}
        disabled={disabled}
        aria-label={`${delta < 0 ? "Decrease" : "Increase"} by ${Math.abs(delta)}`}
        className={`h-8 min-w-8 rounded-lg border-2 px-1.5 flex items-center justify-center font-semibold text-sm transition-colors ${
          disabled
            ? "border-gray-700 bg-gray-800 text-gray-600 cursor-not-allowed"
            : "border-gray-600 bg-gray-700 text-white hover:bg-gray-600 hover:border-gray-500 active:bg-gray-500 cursor-pointer"
        }`}
      >
        {label}
      </button>
    );
  };

  return (
    <div className="flex w-full max-w-xs flex-col gap-2">
      <div className="flex items-center gap-1">
        {coarseStep > 1 && button(-coarseStep, `−${coarseStep}`)}
        {button(-1, "−")}

        <div className="mx-1 h-8 min-w-14 flex-1 rounded-lg border-2 border-gray-600 bg-gray-800 flex items-center justify-center font-semibold text-sm text-white">
          {value}
        </div>

        {button(1, "+")}
        {coarseStep > 1 && button(coarseStep, `+${coarseStep}`)}
      </div>

      {hasSlider && (
        <div>
          {/* Stepped by 1 so the slider can reach every legal number, with the
              buttons either side of it for exact adjustments. */}
          <input
            type="range"
            min={min}
            max={max}
            step={1}
            value={value}
            onChange={(event) => onChange(clamp(Number(event.target.value)))}
            aria-label="Amount"
            className="w-full cursor-pointer accent-blue-500"
          />
          <div className="flex justify-between text-[0.65rem] text-gray-400">
            <span>{min}</span>
            <span>{max}</span>
          </div>
        </div>
      )}
    </div>
  );
}
