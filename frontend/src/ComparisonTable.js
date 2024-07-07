import React from 'react';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from '@mui/material';

function ComparisonTable({ data }) {
  const symbols = Object.keys(data);

  if (symbols.length === 0) {
    return null;
  }

  const metrics = [
    { key: 'price', label: 'Price' },
    { key: 'revenue', label: 'Revenue' },
    { key: 'costOfRevenue', label: 'Cost of Revenue' },
    { key: 'grossProfit', label: 'Gross Profit' },
    { key: 'operatingExpense', label: 'Operating Expense' },
    { key: 'operatingIncome', label: 'Operating Income' },
    { key: 'netNonOperatingInterestIncomeExpense', label: 'Net Non-Operating Interest Income/Expense' },
    { key: 'otherIncomeExpense', label: 'Other Income/Expense' },
    { key: 'pretaxIncome', label: 'Pretax Income' },
    // ... add more metrics as needed from FinancialData
  ];

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Metric</TableCell>
            {symbols.map(symbol => (
              <TableCell key={symbol}>{symbol}</TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {metrics.map(({ key, label }) => (
            <TableRow key={key}>
              <TableCell>{label}</TableCell>
              {symbols.map(symbol => (
                <TableCell key={symbol}>{data[symbol][key]}</TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default ComparisonTable;