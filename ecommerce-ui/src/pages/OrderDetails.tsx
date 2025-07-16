import { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';

interface OrderItem {
  product_id: number;
  product_name: string;
  unit_price: number;
  quantity: number;
}
interface OrderMeta {
  id: number;
  total_amount: number;
  status: string;
  created_at: string;
}

export default function OrderDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [order, setOrder] = useState<OrderMeta | null>(null);
  const [items, setItems] = useState<OrderItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get(`/orders/${id}`)
      .then(res => {
        setOrder(res.data.order);
        setItems(res.data.items);
      })
      .catch(() => setError('Failed to load order details'));
  }, [id]);

  if (error) return <p style={{ color: 'red' }}>{error}</p>;
  if (!order) return <p>Loading order…</p>;

  return (
    <div>
      <button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        ← Back
      </button>
      <h1>Order #{order.id}</h1>
      <p>Status: {order.status}</p>
      <p>Total: ${order.total_amount.toFixed(2)}</p>
      <p>Placed on: {new Date(order.created_at).toLocaleString()}</p>

      <h2>Items</h2>
      <ul>
        {items.map(item => (
          <li key={item.product_id}>
            {item.product_name} × {item.quantity} — $
            {(item.unit_price * item.quantity).toFixed(2)}
          </li>
        ))}
      </ul>
    </div>
  );
}
