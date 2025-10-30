import React from 'react';
import { CircularProgress } from '@mui/material';

export const ButtonSpinner: React.FC<{ size?: number; color?: 'inherit' | 'primary' | 'secondary' | 'error' }> = ({ size = 20, color = 'inherit' }) => (
  <CircularProgress size={size} color={color as any} />
);

export default ButtonSpinner;
