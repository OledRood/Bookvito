import { useEffect, useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export interface LocationItem {
  id: string;
  name: string;
  address?: string;
}

export const useLocationsList = () => {
  const [locations, setLocations] = useState<LocationItem[] | null>(null);
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  useEffect(() => {
    let mounted = true;
    const fetch = async () => {
      setLoading(true);
      try {
        const resp = await api.get<LocationItem[]>('locations/getAll');
        if (!mounted) return;
        setLocations(resp.data || []);
      } catch (err: any) {
        const msg = err?.response?.data?.message || err.message || 'Ошибка при загрузке локаций';
        showNotification(msg, 'error');
        setLocations([]);
      } finally {
        if (mounted) setLoading(false);
      }
    };
    fetch();
    return () => { mounted = false; };
  }, [showNotification]);

  return { locations, loading };
};

export default useLocationsList;
