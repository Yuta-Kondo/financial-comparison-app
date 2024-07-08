import React, { useState } from 'react';
import { Container, Typography, TextField, Button, Box, Alert } from '@mui/material';
import ComparisonTable from './ComparisonTable';
import axios from 'axios';

const REACT_APP_API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [symbols, setSymbols] = useState([]);
  const [inputSymbol, setInputSymbol] = useState('');
  const [comparisonData, setComparisonData] = useState({});
  const [error, setError] = useState(null);

  const handleAddSymbol = async () => {
    console.log("button clicked");
    console.log("input symbol: ", inputSymbol);
    if (inputSymbol && !symbols.includes(inputSymbol.toUpperCase())) {
      const newSymbols = [...symbols, inputSymbol.toUpperCase()];
      setSymbols(newSymbols);
      setInputSymbol('');
      await fetchComparisonData(newSymbols);
      console.log("new symbol: ", newSymbols);
    }
  };

  const fetchComparisonData = async (symbolsToFetch) => {
    try {
      console.log("fetching comparison data");
      setError(null);
      // const response = await axios.get(`${REACT_APP_API_URL}/compare`, {
      //   params: { symbols: symbolsToFetch },
      // });
      const response = await axios.get(`https://financial-comparison-backend.onrender.com/compare`, {
        params: { symbols: symbolsToFetch },
      });
      console.log("response: ", response);
      setComparisonData(response.data);
    } catch (error) {
      console.error('Error fetching comparison data:', error);
      setError('Failed to fetch data. Please make sure the backend server is running.');
    }
  };

  return (
    <Container maxWidth="md">
      <Typography variant="h2" sx={{ mb: 4 }}>
        Financials Comparison
      </Typography>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 4 }}>
        <TextField
          label="Enter stock symbol"
          value={inputSymbol}
          onChange={(e) => setInputSymbol(e.target.value)}
          sx={{ mr: 2 }}
        />
        <Button variant="contained" color="primary" onClick={handleAddSymbol}>
          Add
        </Button>
      </Box>
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      <ComparisonTable data={comparisonData} />
    </Container>
  );
}

export default App;