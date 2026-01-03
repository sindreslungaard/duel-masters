interface ButtonProps {
  disabled?: boolean;
  disabledTooltip?: string;
  onClick?: () => void;
  children: React.ReactNode;
  variant?: "default" | "destructive" | "outline" | "gray";
}

export function Button({
  disabled = false,
  disabledTooltip,
  onClick,
  children,
  variant = "default",
}: ButtonProps) {
  const getColorClasses = () => {
    if (disabled) {
      return "from-gray-600 to-gray-700 cursor-not-allowed opacity-60";
    }
    
    if (variant === "destructive") {
      return "from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 hover:shadow-xl cursor-pointer active:scale-95";
    }
    
    if (variant === "outline") {
      return "bg-transparent border-2 border-white text-white hover:border-gray-300 hover:shadow-xl cursor-pointer active:scale-95";
    }
    
    if (variant === "gray") {
      return "from-gray-500 to-gray-600 hover:from-gray-600 hover:to-gray-700 hover:shadow-xl cursor-pointer active:scale-95";
    }
    
    return "from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 hover:shadow-xl cursor-pointer active:scale-95";
  };

  return (
    <div className="group/button relative">
      <button
        disabled={disabled}
        onClick={onClick}
        className={`w-full px-4 py-2 ${variant === "outline" ? "" : "bg-gradient-to-r text-white"} text-xs font-medium rounded-md shadow-lg transition-all duration-200 whitespace-nowrap overflow-hidden text-ellipsis ${getColorClasses()}`}
      >
        {children}
      </button>
      {disabled && disabledTooltip && (
        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded whitespace-nowrap opacity-0 group-hover/button:opacity-100 transition-opacity pointer-events-none">
          {disabledTooltip}
        </div>
      )}
    </div>
  );
}
