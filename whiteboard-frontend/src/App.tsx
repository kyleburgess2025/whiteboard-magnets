import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Layout from "./pages/Layout";
import Home from "./pages/Home";
import Whiteboard from "./pages/Whiteboard";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="whiteboard" element={<Whiteboard />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
