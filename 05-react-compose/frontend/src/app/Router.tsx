import { BrowserRouter, Routes, Route } from "react-router";
import App from "./App.tsx";
import Products from "@/pages/Products.tsx";

function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<App />}>
          <Route path="/" element={<Products />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default Router;
