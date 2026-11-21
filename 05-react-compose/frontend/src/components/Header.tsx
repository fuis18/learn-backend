import { Button } from "@/components/ui/button";
import { useThemeStore } from "@/store/themeStore";
import { Moon, Sun } from "lucide-react";

const Header = () => {
  const { theme, toggleTheme } = useThemeStore();
  return (
    <header className="flex justify-end py-2 px-4">
      <Button
        variant="outline"
        size="icon"
        onClick={toggleTheme}
        aria-label={
          theme === "dark" ? "Switch to light mode" : "Switch to dark mode"
        }
      >
        {theme === "dark" ? (
          <Sun className="h-4 w-4" />
        ) : (
          <Moon className="h-4 w-4" />
        )}
      </Button>
    </header>
  );
};

export default Header;
