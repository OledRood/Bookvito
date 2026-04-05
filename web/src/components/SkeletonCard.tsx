import React from 'react';
import { Box, Skeleton } from '@mui/material';

const SkeletonCard: React.FC = () => {
  return (
    <Box sx={{ borderRadius: '12px', overflow: 'hidden', boxShadow: 'var(--shadow-dp1)', bgcolor: 'background.paper', height: '100%' }}>
      <Skeleton variant="rectangular" animation="wave" sx={{ height: { xs: 176, sm: 184, md: 192 } }} />
      <Box sx={{ p: 1.5 }}>
        <Skeleton variant="text" animation="wave" width="80%" height={28} />
        <Skeleton variant="text" animation="wave" width="60%" height={20} />
      </Box>
    </Box>
  );
};

export default SkeletonCard;
