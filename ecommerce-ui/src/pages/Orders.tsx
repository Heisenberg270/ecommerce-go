import { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

interface Order {
  id: number;
  total_amount: number;
  status: string;
  created_at: string;
}

export default function Orders() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get('/orders')
      .then(res => setOrders(res.data))
      .catch(() => setError('Failed to fetch orders'));
  }, []);

  if (error) return <p style={{ color: 'red' }}>{error}</p>;

  return (
    <div>
      <h1>Your Orders</h1>
      <ul>
        {orders.map(o => (
          <li key={o.id}>
            <Link to={`/orders/${o.id}`}>Order #{o.id}</Link> — $
            {o.total_amount.toFixed(2)} —{' '}
            {new Date(o.created_at).toLocaleString()}
          </li>
        ))}
      </ul>
    </div>
  );
}
