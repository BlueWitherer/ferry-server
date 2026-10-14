import "./Log.mts";

import './index.css';

import { Box, Paper, Typography, Button, IconButton } from '@mui/material';

import { SiGeode, SiGithub, SiKofi, SiPatreon, SiX, SiYoutube } from "react-icons/si";

function App() {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                p: 2
            }}
        >
            <Paper
                elevation={6}
                sx={{
                    p: 5,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    textAlign: 'center',
                    bgcolor: 'rgba(0, 0, 0, 0.6)',
                    backdropFilter: 'blur(10px)',
                    borderRadius: 4,
                    maxWidth: 500,
                    width: '100%',
                    color: 'white'
                }}
            >
                <Box sx={{ mb: 3 }}>
                    <a href='https://geode-sdk.org/mods/cheeseworks.ferry' target="_blank" aria-label="Download Geode Mod">
                        <img src='https://api.geode-sdk.org/v1/mods/cheeseworks.ferry/logo' className="logo" alt="Ferry logo" style={{ height: '8em' }} />
                    </a>
                </Box>

                <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
                    Ferry
                </Typography>

                <Typography variant="body1" sx={{ mb: 4, color: 'rgba(255, 255, 255, 0.7)' }}>
                    Back up your game client settings!
                </Typography>

                <Button
                    variant="contained"
                    size="large"
                    startIcon={<SiGeode />}
                    href='https://geode-sdk.org/mods/cheeseworks.ferry'
                    target="_blank"
                    sx={{
                        borderRadius: '50px',
                        px: 5,
                        py: 1.5,
                        fontSize: '1.2rem',
                        textTransform: 'none',
                        backgroundColor: '#24292e',
                        '&:hover': {
                            backgroundColor: '#444c56'
                        }
                    }}
                >
                    Download Now
                </Button>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', mt: 2.5 }}>
                    <IconButton
                        className="social-button"
                        component="a"
                        href="https://www.youtube.com/@cheese_works/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="YouTube"
                    >
                        <SiYoutube />
                    </IconButton>

                    <IconButton
                        className="social-button"
                        component="a"
                        href="https://www.github.com/BlueWitherer/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="GitHub"
                    >
                        <SiGithub />
                    </IconButton>

                    <IconButton
                        className="social-button"
                        component="a"
                        href="https://kofi.cheeseworks.gay/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="Kofi"
                    >
                        <SiKofi />
                    </IconButton>

                    <IconButton
                        className="social-button"
                        component="a"
                        href="https://patreon.cheeseworks.gay/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="Patreon"
                    >
                        <SiPatreon />
                    </IconButton>

                    <IconButton
                        className="social-button"
                        component="a"
                        href="https://www.x.com/chris_rhatt/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="X"
                    >
                        <SiX />
                    </IconButton>
                </Box>
            </Paper>
        </Box>
    );
};

export default App;