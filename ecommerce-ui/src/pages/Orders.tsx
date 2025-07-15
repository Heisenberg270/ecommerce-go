import { useEffect, useState } from "react";
import api from "../api";
import { Link } from "react-router-dom";

export default function Orders() {
  const [orders, setOrders] = useState<Array<{
    id: number;
    total_amount: number;
    status: string;
    created_at: string;
  }>>([]);
  const [error, setError] = useState<string>();

  // fetch list of orders
  useEffect(() => {
    api.get("/orders")
      .then(res => setOrders(res.data))
      .catch(() => setError("Failed to load orders"));
  }, []);

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (orders.length === 0) return <p>No orders yet.</p>;

  return (
    <div>
      <h1>My Orders</h1>
      <ul>
       {orders.map(o => (
         <li key={o.id}>
           <Link to={`/orders/${o.id}`}>
             Order #{o.id}: ${o.total_amount.toFixed(2)} — {o.status} on{" "}
             {new Date(o.created_at).toLocaleString()}
           </Link>
         </li>
       ))}
      </ul>
    </div>
  );
}
