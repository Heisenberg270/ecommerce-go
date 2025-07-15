import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import api from "../api";

type OrderInfo = {
  order: {
    id: number;
    user_id: number;
    total_amount: number;
    status: string;
    created_at: string;
  };
  items: Array<{
    product_id: number;
    quantity: number;
    unit_price: number;
    product_name: string;
  }>;
};

export default function OrderDetails() {
  const { orderID } = useParams<{ orderID: string }>();
  const [info, setInfo] = useState<OrderInfo | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!orderID) return;
    api
      .get<OrderInfo>(`/orders/${orderID}`)
      .then(res => setInfo(res.data))
      .catch(() => setError("Failed to load order"));
  }, [orderID]);

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (!info) return <p>Loading order…</p>;

  return (
    <div>
      <h1>Order #{info.order.id}</h1>
      <p>
        <strong>Status:</strong> {info.order.status}<br/>
        <strong>Total:</strong> ${info.order.total_amount.toFixed(2)}<br/>
        <strong>Placed at:</strong>{" "}
        {new Date(info.order.created_at).toLocaleString()}
      </p>

      <h2>Items</h2>
      {info.items.length === 0 ? (
        <p>(No items in this order.)</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Product</th><th>Qty</th><th>Unit Price</th><th>Subtotal</th>
            </tr>
          </thead>
          <tbody>
            {info.items.map(it => (
              <tr key={it.product_id}>
                <td>{it.product_name}</td>
                <td>{it.quantity}</td>
                <td>${it.unit_price.toFixed(2)}</td>
                <td>${(it.quantity * it.unit_price).toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <p>
        <Link to="/orders">← Back to Orders</Link>
      </p>
    </div>
  );
}
