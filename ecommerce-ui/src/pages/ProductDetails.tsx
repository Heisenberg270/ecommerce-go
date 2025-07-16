import { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface Product { id: number; name: string; description: string; price: number; }

export default function ProductDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [product, setProduct] = useState<Product | null>(null);
  const [qty, setQty] = useState(1);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get(`/products/${id}`)
      .then(res => setProduct(res.data))
      .catch(() => setError('Product not found'));
  }, [id]);

  const addToCart = async () => {
    if (!isAuthenticated) return navigate('/login');
    try {
      let cartID = localStorage.getItem('cartID');
      if (!cartID) {
        const resp = await axios.post('/carts');
        cartID = String(resp.data.id);
        localStorage.setItem('cartID', cartID);
      }
      await axios.post(`/carts/${cartID}/items`, {
        product_id: Number(id),
        quantity: qty,
      });
      navigate('/cart');
    } catch {
      alert('Failed to add to cart');
    }
  };

  if (error) return <p style={{ color: 'red' }}>{error}</p>;
  if (!product) return <p>Loading product…</p>;

  return (
    <div>
      <h1>{product.name}</h1>
      <p>{product.description}</p>
      <strong>${product.price.toFixed(2)}</strong>
      <div style={{ marginTop: 12 }}>
        Qty:{' '}
        <input
          type="number"
          min={1}
          value={qty}
          onChange={e => setQty(Math.max(1, +e.target.value))}
          style={{ width: 60 }}
        />
        <button onClick={addToCart} style={{ marginLeft: 8 }}>
          Add to Cart
        </button>
      </div>
    </div>
  );
}

