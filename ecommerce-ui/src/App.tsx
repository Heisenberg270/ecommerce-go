import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import Products from "./pages/Products";
import Signup from "./pages/Signup";
import Login from "./pages/Login";
import Cart from "./pages/Cart";
import Orders from "./pages/Orders";
import OrderDetails from "./pages/OrderDetails";

function App() {
  return (
    <BrowserRouter>
      <nav style={{ padding: "1rem", borderBottom: "1px solid #ccc" }}>
        <Link to="/" style={{ marginRight: "1rem" }}>Products</Link>
        <Link to="/signup" style={{ marginRight: "1rem" }}>Signup</Link>
        <Link to="/login" style={{ marginRight: "1rem" }}>Login</Link>
        <Link to="/cart" style={{ marginRight: "1rem" }}>Cart</Link>
        <Link to="/orders" style={{ marginLeft: "1rem" }}>Orders</Link>

      </nav>
      <main style={{ padding: "1rem" }}>
        <Routes>
          <Route path="/" element={<Products />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/login" element={<Login />} />
          <Route path="/cart" element={<Cart />} />
          <Route path="/orders" element={<Orders />} />
          <Route path="/orders/:orderID" element={<OrderDetails />} />
          {/* future: protected /cart & /orders */}
        </Routes>
      </main>
    </BrowserRouter>
  );
}

export default App;
