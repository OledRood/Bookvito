import { useEffect, useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export interface LocationItem {
  id: string;
  name: string;
  address?: string;
}

let locationsCache: LocationItem[] | null = null;
let locationsPromise: Promise<LocationItem[]> | null = null;

const fetchLocations = async (): Promise<LocationItem[]> => {
  if (locationsCache) return locationsCache;
  if (locationsPromise) return locationsPromise;

  locationsPromise = api.get<LocationItem[]>('locations/getAll')
    .then((resp) => {
      const normalized = Array.isArray(resp.data) ? resp.data : [];
      locationsCache = normalized;
      return normalized;
    })
    .finally(() => {
      locationsPromise = null;
    });

  return locationsPromise;
};

export const useLocationsList = () => {
  const [locations, setLocations] = useState<LocationItem[] | null>(null);
  const [loading, setLoading] = useState(!locationsCache);
  const { showNotification } = useNotification();

  useEffect(() => {
    let mounted = true;

    if (locationsCache) {
      setLocations(locationsCache);
      setLoading(false);
      return () => { mounted = false; };
    }

    const fetch = async () => {
      setLoading(true);
      try {
        const resp = await fetchLocations();
        if (!mounted) return;
        setLocations(resp);
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
