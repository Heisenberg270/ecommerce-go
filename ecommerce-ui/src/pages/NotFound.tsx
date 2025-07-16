import { Link } from 'react-router-dom';

export default function NotFound() {
  return (
    <div style={{ textAlign: 'center', marginTop: 80 }}>
      <h1>404: Page Not Found</h1>
      <p>
        Oops! That page doesn’t exist. <Link to="/">Go back home</Link>.
      </p>
    </div>
  );
}
