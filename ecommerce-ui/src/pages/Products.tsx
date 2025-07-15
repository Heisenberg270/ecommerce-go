// src/pages/Products.tsx
import { useEffect, useState } from "react";
import api from "../api";
import type { Product } from "../api";

export default function Products() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string|undefined>(undefined);

  useEffect(() => {
    api.get<Product[]>("/products")
      .then(res => setProducts(res.data))
      .catch(err => {
        console.error(err);
        setError("Failed to load products");
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading products…</p>;
  if (error)   return <p style={{ color: "red" }}>{error}</p>;

  return (
    <div>
      <h1>Products</h1>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {products.map(p => (
          <li key={p.id} style={{ marginBottom: "1rem" }}>
            <strong>{p.name}</strong> — ${p.price.toFixed(2)}
            <p>{p.description}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}
