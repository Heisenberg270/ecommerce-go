// src/pages/Cart.tsx
import { useEffect, useState } from "react";
import api from "../api";

type CartInfo = {
  cart: {
    id: number;
    user_id: number;
    created_at: string;
  };
  items: Array<{
    cart_id: number;
    product_id: number;
    quantity: number;
    product_name: string;
    unit_price: number;
  }>;
};

export default function Cart() {
  const [cartID, setCartID] = useState<number | null>(null);
  const [info, setInfo]     = useState<CartInfo | null>(null);
  const [error, setError]   = useState<string | null>(null);

  const [prodID, setProdID]     = useState<number>(0);
  const [quantity, setQuantity] = useState<number>(1);

  // 1) Initialize or fetch cartID
  useEffect(() => {
    const existing = localStorage.getItem("cartID");
    async function init() {
      try {
        let id = existing ? Number(existing) : null;
        if (!id) {
          const resp = await api.post<{ id: number }>("/carts", {});
          id = resp.data.id;
          localStorage.setItem("cartID", id.toString());
        }
        setCartID(id);
      } catch {
        setError("Could not initialize cart");
      }
    }
    init();
  }, []);

  // 2) Fetch cart contents whenever cartID is known
  useEffect(() => {
    if (cartID == null) return;
    async function fetchCart() {
      try {
        const resp = await api.get<CartInfo>(`/carts/${cartID}`);
        setInfo(resp.data);
      } catch {
        setError("Failed to load cart");
      }
    }
    fetchCart();
  }, [cartID]);

  // 3) Add item
  const addItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (cartID == null) return;
    try {
      await api.post(`/carts/${cartID}/items`, { product_id: prodID, quantity });
      const resp = await api.get<CartInfo>(`/carts/${cartID}`);
      setInfo(resp.data);
    } catch {
      setError("Add item failed");
    }
  };

  // 4) Remove item
  const removeItem = async (productID: number) => {
    if (cartID == null) return;
    try {
      await api.delete(`/carts/${cartID}/items/${productID}`);
      const resp = await api.get<CartInfo>(`/carts/${cartID}`);
      setInfo(resp.data);
    } catch {
      setError("Remove failed");
    }
  };

  // —— EARLY RETURNS —— //
  if (error) {
    return <p style={{ color: "red" }}>{error}</p>;
  }
  if (cartID == null) {
    // Still determining cartID
    return <p>Initializing cart…</p>;
  }
  if (info == null) {
    // cartID known, but info not yet fetched
    return <p>Loading cart…</p>;
  }

  // Normalize items array so it's never null
  const items = info.items || [];

  // —— MAIN RENDER —— //
  return (
    <div>
      <h1>My Cart (#{info.cart.id})</h1>

      <form onSubmit={addItem} style={{ marginBottom: "1rem" }}>
        <input
          type="number"
          placeholder="Product ID"
          value={prodID}
          min={1}
          onChange={(e) => setProdID(Number(e.target.value))}
          required
          style={{ marginRight: ".5rem" }}
        />
        <input
          type="number"
          placeholder="Quantity"
          value={quantity}
          min={1}
          onChange={(e) => setQuantity(Number(e.target.value))}
          required
          style={{ marginRight: ".5rem" }}
        />
        <button type="submit">Add Item</button>
      </form>

      {items.length === 0 ? (
        <p>Your cart is empty.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Product</th>
              <th>Qty</th>
              <th>Unit Price</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map((it) => (
              <tr key={it.product_id}>
                <td>{it.product_name}</td>
                <td>{it.quantity}</td>
                <td>${it.unit_price.toFixed(2)}</td>
                <td>
                  <button onClick={() => removeItem(it.product_id)}>
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>

          <button
            onClick={async () => {
              if (!cartID) return;
              try {
                await api.post("/orders", { cart_id: cartID });
                // clear local cart ID so a new one starts fresh
                localStorage.removeItem("cartID");
                // redirect to Orders page
                window.location.href = "/orders";
              } catch {
                setError("Checkout failed");
              }
            }}
          >
            Place Order
          </button>

        </table>

        
      )}
    </div>
  );
}
