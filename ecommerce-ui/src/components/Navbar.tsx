import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { isAuthenticated, logout } = useAuth();

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography
          variant="h6"
          component={RouterLink}
          to="/"
          sx={{ flexGrow: 1, textDecoration: 'none', color: 'inherit' }}
        >
          MyStore
        </Typography>
        <Button component={RouterLink} to="/products" color="inherit">
          Products
        </Button>
        {isAuthenticated ? (
          <>
            <Button component={RouterLink} to="/cart" color="inherit">
              Cart
            </Button>
            <Button component={RouterLink} to="/orders" color="inherit">
              Orders
            </Button>
            <Button onClick={logout} color="inherit">
              Logout
            </Button>
          </>
        ) : (
          <>
            <Button component={RouterLink} to="/login" color="inherit">
              Login
            </Button>
            <Button component={RouterLink} to="/signup" color="inherit">
              Signup
            </Button>
          </>
        )}
      </Toolbar>
    </AppBar>
  );
}
