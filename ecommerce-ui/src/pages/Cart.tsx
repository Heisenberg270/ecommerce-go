// src/pages/Cart.tsx
import { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

interface CartItem {
  product_id: number;
  product_name: string;
  unit_price: number;
  quantity: number;
}

export default function Cart() {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // Redirect away if not logged in
  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

  const [cartID, setCartID] = useState<number | null>(null);
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Initialize or create cart on mount
  useEffect(() => {
    if (!isAuthenticated) return;

    const stored = localStorage.getItem('cartID');
    if (stored) {
      const id = Number(stored);
      setCartID(id);
      fetchCart(id);
    } else {
      axios
        .post('/carts')
        .then(res => {
          const newID = res.data.id as number;
          localStorage.setItem('cartID', String(newID));
          setCartID(newID);
          fetchCart(newID);
        })
        .catch(() => {
          setError('Failed to create cart');
          setLoading(false);
        });
    }
  }, [isAuthenticated]);

  /**
   * Fetches cart items (silent = no loading flash).
   */
  const fetchCart = (id: number, silent = false) => {
    if (!silent) setLoading(true);
    axios
      .get(`/carts/${id}`)
      .then(res => {
        // always sort by product_id for stable ordering
        const arr = Array.isArray(res.data.items) ? res.data.items : [];
        arr.sort((a: any, b: any) => a.product_id - b.product_id);
        setItems(arr);
        if (!silent) setLoading(false);
      })
      .catch(() => {
        setError('Failed to load cart');
        if (!silent) setLoading(false);
      });
  };

  /**
   * Update quantity by deleting the old line then re-adding with newQty.
   */
  const updateQuantity = async (productId: number, newQty: number) => {
    if (!cartID) return;
    try {
      // 1️⃣ remove whatever is there
      await axios.delete(`/carts/${cartID}/items/${productId}`);
      // 2️⃣ if they still want >0, re-insert with that full qty
      if (newQty > 0) {
        await axios.post(`/carts/${cartID}/items`, {
          product_id: productId,
          quantity: newQty,
        });
      }
      // 3️⃣ silent refresh so we don't flash "Loading…"
      fetchCart(cartID, true);
    } catch (e) {
      console.error(e);
      alert('Failed to update quantity');
    }
  };

  /**
   * Remove an item completely.
   */
  const removeItem = async (productId: number) => {
    if (!cartID) return;
    try {
      await axios.delete(`/carts/${cartID}/items/${productId}`);
      fetchCart(cartID, true);
    } catch {
      alert('Failed to remove item');
    }
  };

  /**
   * Place the order (checkout).
   */
 const handleCheckout = async () => {
  if (!cartID) return;
  try {
    const { data: order } = await axios.post('/orders', { cart_id: cartID });
    localStorage.removeItem('cartID');
    navigate(`/orders/${order.id}`);
  } catch (err: any) {
    if (err.response?.status === 401) navigate('/login');
    else alert('Failed to place order');
  }
};


  // —— RENDERING ——

  if (!isAuthenticated) {
    // Redirect effect will run; render nothing temporarily
    return null;
  }

  if (loading) {
    return <p>Loading your cart…</p>;
  }

  if (error) {
    return <p style={{ color: 'red' }}>{error}</p>;
  }

  if (items.length === 0) {
    return <p>Your cart is empty.</p>;
  }

  return (
    <div>
      <h1>Your Cart</h1>
      <ul>
        {items.map(item => (
          <li key={item.product_id} style={{ marginBottom: 12 }}>
            <strong>{item.product_name}</strong> — ${item.unit_price.toFixed(2)}
            <div style={{ marginTop: 4 }}>
              Qty:{' '}
              <input
                type="number"
                min={0}
                value={item.quantity}
                onChange={e =>
                  updateQuantity(
                    item.product_id,
                    Math.max(0, +e.target.value)
                  )
                }
                style={{ width: 60 }}
              />
              <button
                onClick={() => removeItem(item.product_id)}
                style={{ marginLeft: 8 }}
              >
                Remove
              </button>
            </div>
          </li>
        ))}
      </ul>
      <button onClick={handleCheckout}>Place Order</button>
    </div>
  );
}
