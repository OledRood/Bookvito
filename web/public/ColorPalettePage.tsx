// TO_BE_REMOVED_BEFORE_RELEASE
import React from 'react';
import { Box, Typography, Paper, Grid, useTheme, Container } from '@mui/material';
import { useNotification } from './NotificationContext';
import {
    alpha,
    darken,
    lighten
} from '@mui/material/styles';

const ColorBox: React.FC<{ color: string; name: string; path?: string }> = ({ color, name, path }) => {
    const theme = useTheme();
    const { showNotification } = useNotification();
    let textColor;
    try {
        textColor = theme.palette.getContrastText(color);
    } catch (e) {
        textColor = theme.palette.text.primary;
    }

    const valueToCopy = path ?? color;

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(valueToCopy);
            showNotification(`Скопировано ${valueToCopy}`, 'success');
        } catch (err) {
            // Fallback: select temporary input
            try {
                const input = document.createElement('input');
                input.value = valueToCopy;
                document.body.appendChild(input);
                input.select();
                document.execCommand('copy');
                document.body.removeChild(input);
                showNotification(`Скопировано ${valueToCopy}`, 'success');
            } catch (e) {
                showNotification('Не удалось скопировать значение', 'error');
            }
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleCopy();
        }
    };

    return (
        <Paper
            elevation={2}
            role="button"
            tabIndex={0}
            onClick={handleCopy}
            onKeyDown={handleKeyDown}
            aria-label={`Копировать ${path ?? name} ${color}`}
            sx={{
                backgroundColor: color,
                color: textColor,
                p: 2,
                borderRadius: 2,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                height: 100,
                cursor: 'pointer',
                userSelect: 'none',
            }}
        >
            <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                {name}
            </Typography>
            <Typography variant="caption" sx={{ wordBreak: 'break-all' }}>
                {color}
            </Typography>
        </Paper>
    );
};

const PaletteSection: React.FC<{ title: string; colors: Record<string, string> }> = ({ title, colors }) => {
    return (
        <Box sx={{ mb: 4 }}>
            <Typography variant="h5" component="h2" gutterBottom>
                {title}
            </Typography>
            <Box
                sx={{
                    display: 'grid',
                    gap: 2,
                    gridTemplateColumns: {
                        xs: 'repeat(2, 1fr)',
                        sm: 'repeat(3, 1fr)',
                        md: 'repeat(4, 1fr)',
                        lg: 'repeat(6, 1fr)'
                    }
                }}
            >
                {Object.entries(colors).map(([name, color]) =>
                    typeof color === 'string' ? (
                        <Box key={name}>
                            <ColorBox name={name} color={color} path={`${title.toLowerCase()}.${name}`} />
                        </Box>
                    ) : null,
                )}
            </Box>
        </Box>
    );
};

const ColorPalettePage: React.FC = () => {
    const theme = useTheme();
    const { palette } = theme;

    // Collect CSS custom properties that start with `--md-sys-color-` (Material You / design tokens)
    // so we can display colors like `surface-container` which are defined as CSS vars.
    const getCssSystemColors = (): Record<string, string> => {
        if (typeof window === 'undefined') return {};
        const styles = getComputedStyle(document.documentElement);
        const res: Record<string, string> = {};
        for (let i = 0; i < styles.length; i++) {
            const prop = styles.item(i);
            if (!prop) continue;
            if (prop.startsWith('--md-sys-color-')) {
                const value = styles.getPropertyValue(prop).trim();
                const key = prop.replace('--md-sys-color-', '');
                res[key] = value;
            }
        }
        return res;
    };

    const paletteSections = [
        { title: 'Primary', colors: palette.primary as unknown as Record<string, string> },
        { title: 'Secondary', colors: palette.secondary as unknown as Record<string, string> },
        { title: 'Error', colors: palette.error as unknown as Record<string, string> },
        { title: 'Warning', colors: palette.warning as unknown as Record<string, string> },
        { title: 'Info', colors: palette.info as unknown as Record<string, string> },
        { title: 'Success', colors: palette.success as unknown as Record<string, string> },
        { title: 'Text', colors: palette.text as unknown as Record<string, string> },
        { title: 'Background', colors: palette.background as unknown as Record<string, string> },
        { title: 'Action', colors: palette.action as unknown as Record<string, string> },
        { title: 'Common', colors: palette.common as unknown as Record<string, string> },
        { title: 'System CSS vars', colors: getCssSystemColors() },
    ];

    return (
        <Container maxWidth="lg">
            <Typography variant="h4" component="h1" gutterBottom sx={{ mt: 4, mb: 4 }}>
                Палетка цветов темы
            </Typography>
            {paletteSections.map(
                (section) =>
                    section.colors && <PaletteSection key={section.title} title={section.title} colors={section.colors} />,
            )}
        </Container>
    );
};

export default ColorPalettePage;