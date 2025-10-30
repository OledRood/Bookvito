import React from 'react';
import { Box, Skeleton } from '@mui/material';

const SkeletonCard: React.FC = () => {
  return (
    <Box sx={{ borderRadius: '12px', overflow: 'hidden', boxShadow: 'var(--shadow-dp1)', bgcolor: 'background.paper' }}>
      <Skeleton variant="rectangular" animation="wave" height={190} />
      <Box sx={{ p: 2 }}>
        <Skeleton variant="text" animation="wave" width="80%" height={28} />
        <Skeleton variant="text" animation="wave" width="60%" height={20} />
      </Box>
    </Box>
  );
};

export default SkeletonCard;
