import React from 'react';
import { Container, Box } from '@mui/material';
import Navbar from './Navbar';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <Box display="flex" flexDirection="column" minHeight="100vh">
      <Navbar />
      <Container component="main" sx={{ py: 4, flexGrow: 1 }}>
        {children}
      </Container>
      <Box component="footer" py={2} bgcolor="grey.100" textAlign="center">
        © 2025 MyStore
      </Box>
    </Box>
  );
}
