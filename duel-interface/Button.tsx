interface ButtonProps {
  disabled?: boolean;
  disabledTooltip?: string;
  onClick?: () => void;
  children: React.ReactNode;
}

export function Button({
  disabled = false,
  disabledTooltip = "This action is not available right now",
  onClick,
  children,
}: ButtonProps) {
  return (
    <div className="bg-black/50 p-3 rounded-md">
      <div className="group/button relative">
        <button
          disabled={disabled}
          onClick={onClick}
          className={`w-full px-4 py-2 bg-gradient-to-r text-white text-sm font-medium rounded-md shadow-lg transition-all duration-200 ${
            disabled
              ? "from-gray-600 to-gray-700 cursor-not-allowed opacity-60"
              : "from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 hover:shadow-xl cursor-pointer active:scale-95"
          }`}
        >
          {children}
        </button>
        {disabled && disabledTooltip && (
          <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded whitespace-nowrap opacity-0 group-hover/button:opacity-100 transition-opacity pointer-events-none">
            {disabledTooltip}
          </div>
        )}
      </div>
    </div>
  );
}
