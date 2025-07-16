import { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

interface Product { id: number; name: string; description?: string; price: number; }

export default function Products() {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [quantities, setQuantities] = useState<Record<number, number>>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get('/products')
      .then(res => setProducts(res.data))
      .catch(() => setError('Failed to load products'));
  }, []);

  const handleAddToCart = async (productId: number) => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    try {
      // Create or reuse cart
      let cartID = localStorage.getItem('cartID');
      if (!cartID) {
        const resp = await axios.post('/carts');
        cartID = String(resp.data.id);
        localStorage.setItem('cartID', cartID);
      }
      const qty = quantities[productId] || 1;
      await axios.post(`/carts/${cartID}/items`, {
        product_id: productId,
        quantity: qty,
      });
      alert(`Added ${qty}× item to cart`);
    } catch (err: any) {
      // If somehow unauthorized (token expired?)
      if (err.response?.status === 401) {
        navigate('/login');
      } else {
        alert('Failed to add to cart');
      }
    }
  };

  if (error) return <p style={{ color: 'red' }}>{error}</p>;

  return (
    <div>
      <h1>All Products</h1>
      <ul>
        {products.map(p => (
          <li key={p.id} style={{ marginBottom: 16 }}>
            <h2>{p.name}</h2>
            <p>{p.description}</p>
            <strong>${p.price.toFixed(2)}</strong>
            <div style={{ marginTop: 8 }}>
              <label>
                Qty:{' '}
                <input
                  type="number"
                  min={1}
                  value={quantities[p.id] || 1}
                  onChange={e =>
                    setQuantities({
                      ...quantities,
                      [p.id]: Math.max(1, Number(e.target.value)),
                    })
                  }
                  style={{ width: 60 }}
                />
              </label>
              <button
                onClick={() => handleAddToCart(p.id)}
                style={{ marginLeft: 8 }}
              >
                Add to Cart
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
