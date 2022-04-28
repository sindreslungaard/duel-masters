export const toggleMute = player => {
  const settings = getSettings();
  const isMuted = settings.muted.includes(player);

  if (isMuted) {
    settings.muted = settings.muted.filter(p => p !== player);
  } else {
    settings.muted.push(player);
  }

  setSettings(settings);

  return !isMuted;
};

export const isMuted = player => getSettings().muted.includes(player);

export const didSeeMuteWarning = () => {
  const settings = getSettings();

  if (settings.didMute) {
    return true;
  }

  settings.didMute = true;
  setSettings(settings);
  return false;
};

export const getSettings = () =>
  JSON.parse(
    localStorage.getItem("settings") ??
      JSON.stringify({
        muted: [],
        didMute: false,
        noUpsideDownCards: false
      })
  );

export const patchSettings = newSettings => {
  setSettings({ ...getSettings(), ...newSettings });
};

const setSettings = players => {
  localStorage.setItem("settings", JSON.stringify(players));

  // "storage" event does not fire in the same window. fire a custom event
  // setTimeout is needed for vue to have time to reconcile state
  setTimeout(() => dispatchEvent(new Event("storageUpdated")));
};
