import { motion, useScroll, useTransform } from "motion/react";
import { Download, Github, Terminal, Zap, Globe, Layers, Activity, Server } from "lucide-react";

const GITHUB_URL = "https://github.com/raz0red/RADPresence";
const RELEASES_URL = "https://github.com/raz0red/RADPresence/releases";

export default function App() {
  const { scrollYProgress } = useScroll();
  const opacity = useTransform(scrollYProgress, [0, 0.2], [1, 0]);
  const scale = useTransform(scrollYProgress, [0, 0.2], [1, 0.95]);

  return (
    <div className="min-h-screen bg-zinc-950 overflow-x-hidden">
      {/* Abstract Background Elements */}
      <div className="fixed inset-0 z-0 pointer-events-none overflow-hidden">
        <div className="absolute top-[-20%] left-[-10%] w-[50vw] h-[50vw] bg-blurple-500/20 rounded-full blur-[120px] mix-blend-screen opacity-50" />
        <div className="absolute bottom-[-20%] right-[-10%] w-[50vw] h-[50vw] bg-retro-orange/20 rounded-full blur-[120px] mix-blend-screen opacity-50" />
        <div className="absolute top-[40%] left-[60%] w-[30vw] h-[30vw] bg-indigo-500/10 rounded-full blur-[100px] mix-blend-screen opacity-40" />
      </div>

      {/* Nav */}
      <nav className="fixed top-0 w-full z-50 bg-zinc-950/50 backdrop-blur-md border-b border-white/10">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2 font-display font-bold text-xl tracking-tight">
            <span className="text-blurple-500 text-2xl">RAD</span>
            <span className="text-white">Presence</span>
          </div>
          <div className="flex gap-4">
            <a href={GITHUB_URL} target="_blank" rel="noreferrer" className="text-zinc-400 hover:text-white transition-colors flex items-center gap-2 text-sm font-medium">
              <Github size={20} />
              <span className="hidden sm:inline">GitHub</span>
            </a>
            <a href={RELEASES_URL} target="_blank" rel="noreferrer" className="bg-white/10 hover:bg-white/20 text-white px-4 py-2 rounded-full transition-colors flex items-center gap-2 text-sm font-medium">
              <Download size={16} />
              Releases
            </a>
          </div>
        </div>
      </nav>

      <main className="relative z-10 pt-32 pb-24 flex flex-col items-center">
        {/* Split Hero Section */}
        <section className="w-full max-w-6xl mx-auto px-6 flex flex-col lg:flex-row items-center justify-center gap-12 lg:gap-16 mb-32">
          
          <div className="flex-1 text-center lg:text-left flex flex-col items-center lg:items-start pt-8 max-w-2xl">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, ease: "easeOut" }}
              className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blurple-500/10 border border-blurple-500/20 text-blurple-300 text-sm mb-6"
            >
              <Zap size={14} className="text-retro-gold" />
              <span className="font-medium">Cross-platform. Single binary. Zero dependencies.</span>
            </motion.div>

            <motion.h1
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1, ease: "easeOut" }}
              className="font-display font-bold text-4xl sm:text-5xl md:text-6xl lg:text-7xl tracking-tighter leading-[1.05] bg-gradient-to-br from-white via-zinc-200 to-zinc-500 bg-clip-text text-transparent mb-6"
            >
              Show off your <br className="hidden lg:block" />
              <span className="text-blurple-500 bg-clip-text text-transparent bg-gradient-to-r from-blurple-500 to-indigo-400 inline-block break-words max-w-full">RetroAchievements</span>
              <br className="hidden lg:block" /> on Discord.
            </motion.h1>

            <motion.p
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.2, ease: "easeOut" }}
              className="text-lg md:text-xl text-zinc-400 max-w-xl mx-auto lg:mx-0 mb-8 leading-relaxed"
            >
              A silent, native background service that automatically mirrors your RetroAchievements session to your Discord Rich Presence.
            </motion.p>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.3, ease: "easeOut" }}
              className="flex flex-col sm:flex-row gap-4 items-center w-full justify-center lg:justify-start"
            >
              <a
                href={RELEASES_URL}
                target="_blank"
                rel="noreferrer"
                className="group relative inline-flex items-center justify-center gap-2 bg-blurple-500 hover:bg-blurple-600 text-white font-semibold py-3.5 px-8 rounded-full transition-all overflow-hidden shadow-[0_0_40px_-10px_rgba(88,101,242,0.8)] w-full sm:w-auto shrink-0"
              >
                <div className="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out" />
                <Download size={20} className="relative z-10" />
                <span className="relative z-10">Get RAD Presence</span>
              </a>
              <a
                href={GITHUB_URL}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center justify-center gap-2 bg-zinc-900 border border-zinc-700 hover:bg-zinc-800 hover:border-zinc-500 text-zinc-300 font-semibold py-3.5 px-8 rounded-full transition-all w-full sm:w-auto shrink-0"
              >
                <Github size={20} />
                <span>View Source</span>
              </a>
            </motion.div>
          </div>

          <motion.div 
            initial={{ opacity: 0, scale: 0.9, rotate: 2 }}
            animate={{ opacity: 1, scale: 1, rotate: 0 }}
            transition={{ duration: 0.7, delay: 0.2, type: "spring" }}
            className="w-full lg:w-auto flex justify-center relative shrink-0"
          >
            <div className="absolute inset-0 bg-blurple-500/20 blur-[100px] rounded-full scale-75" />
            
            <div className="relative transform-gpu hover:scale-[1.02] transition-transform duration-500 mx-auto shadow-2xl shadow-blurple-500/20 rounded-2xl overflow-hidden border border-white/10 flex justify-center bg-[#1E1F22] w-full max-w-[293px]">
              <img
                src="https://raw.githubusercontent.com/raz0red/RADPresence/master/.github/assets/screenshot.png"
                alt="Discord Rich Presence Showcasing RetroAchievements"
                className="w-full h-auto block"
              />
            </div>
          </motion.div>

        </section>

        {/* Features Bento Grid */}
        <section className="w-full max-w-7xl mx-auto px-6 py-16">
          <div className="text-center mb-16">
            <h2 className="font-display font-bold text-3xl md:text-5xl tracking-tight mb-4 text-white">
              Everything you need,<br />nothing you don't.
            </h2>
            <p className="text-zinc-400 text-lg">Designed to run silently and out of your way.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <FeatureCard 
              icon={<Server size={32} />}
              title="Native Background Service"
              description="Runs natively via Windows SCM, macOS launchd, or Linux systemd. Starts silently on login."
              className="lg:col-span-2 bg-gradient-to-br from-zinc-900 to-zinc-900/50"
            />
            <FeatureCard 
              icon={<Layers size={32} />}
              title="Single Binary"
              description="No runtime needed. No installers. Just drop the executable and run."
              className="bg-gradient-to-br from-zinc-900 to-zinc-900/50"
            />
            <FeatureCard 
              icon={<Globe size={32} />}
              title="Local Web UI"
              description="Built-in, optional local dashboard to view status, tweak settings, and read logs directly from your browser."
              className="bg-gradient-to-br from-zinc-900 to-zinc-900/50"
            />
            <FeatureCard 
              icon={<Activity size={32} />}
              title="Real-time Polling"
              description="Automatically detects when you start and stop playing. Updates Discord with your current progress, cover art, and more."
              className="lg:col-span-2 bg-gradient-to-br from-zinc-900 to-zinc-900/50"
            />
          </div>
        </section>

        {/* Web UI Feature Showcase */}
        <section className="w-full max-w-7xl mx-auto px-6 py-16">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
            <div>
              <h2 className="font-display font-bold text-3xl md:text-5xl tracking-tight mb-6 text-white text-balance">
                Total control,<br />right from your browser.
              </h2>
              <p className="text-zinc-400 text-lg mb-8 leading-relaxed text-balance">
                Enable the optional Web UI to get a live view of your session, 
                tweak settings without restarting the service, and monitor the real-time activity log.
              </p>
              <ul className="space-y-4">
                {[
                  "Live Session Viewer (updates every 3s)",
                  "Real-time colored log tailing",
                  "Modify RA credentials securely",
                  "Toggle visibility of buttons & stats"
                ].map((item, i) => (
                  <li key={i} className="flex items-center gap-3 text-zinc-300">
                    <div className="w-6 h-6 rounded-full bg-blurple-500/20 flex items-center justify-center text-blurple-500 shrink-0">
                      <Zap size={14} />
                    </div>
                    {item}
                  </li>
                ))}
              </ul>
            </div>
            <div className="relative">
              <div className="absolute -inset-4 bg-gradient-to-br from-blurple-500/20 to-indigo-500/20 rounded-3xl blur-2xl z-0" />
              <div className="relative z-10 grid gap-4">
                <motion.img 
                  initial={{ opacity: 0, x: 20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  src="https://raw.githubusercontent.com/raz0red/RADPresence/master/.github/assets/webui-status.png" 
                  alt="Web UI Status" 
                  className="rounded-xl border border-white/10 shadow-2xl skew-y-2 -rotate-2 hover:skew-y-0 hover:rotate-0 transition-transform duration-500 w-full"
                />
                <motion.img 
                  initial={{ opacity: 0, x: 20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true, margin: "-100px" }}
                  transition={{ delay: 0.1 }}
                  src="https://raw.githubusercontent.com/raz0red/RADPresence/master/.github/assets/webui-settings.png" 
                  alt="Web UI Settings" 
                  className="rounded-xl border border-white/10 shadow-2xl -skew-y-2 rotate-2 hover:skew-y-0 hover:rotate-0 transition-transform duration-500 w-3/4 ml-auto -mt-12"
                />
              </div>
            </div>
          </div>
        </section>

        {/* Getting Started Terminal */}
        <section className="w-full max-w-4xl mx-auto px-6 py-32 text-center">
          <h2 className="font-display font-bold text-3xl md:text-5xl tracking-tight mb-8 text-white">
            Get up and running in seconds.
          </h2>
          <div className="bg-[#1C1C1E] border border-white/10 rounded-2xl overflow-hidden shadow-2xl text-left max-w-2xl mx-auto">
            <div className="flex items-center px-4 py-3 bg-[#2D2D2F] border-b border-white/5">
              <div className="flex gap-2">
                <div className="w-3 h-3 rounded-full bg-red-500/80" />
                <div className="w-3 h-3 rounded-full bg-amber-500/80" />
                <div className="w-3 h-3 rounded-full bg-green-500/80" />
              </div>
              <div className="ml-auto flex items-center justify-center font-mono text-[10px] text-zinc-500">
                bash — RADPresence
              </div>
            </div>
            <div className="p-6 font-mono text-sm leading-relaxed overflow-x-auto">
              <div className="flex gap-4">
                <span className="text-zinc-600 select-none">1</span>
                <span><span className="text-blurple-400"># Set your credentials</span></span>
              </div>
              <div className="flex gap-4">
                <span className="text-zinc-600 select-none">2</span>
                <span className="text-white">radpresence set --username <span className="text-retro-orange">YOUR_USER</span> --apikey <span className="text-retro-orange">YOUR_KEY</span></span>
              </div>
              <div className="flex gap-4 mt-4">
                <span className="text-zinc-600 select-none">3</span>
                <span><span className="text-blurple-400"># Install & start the background service</span></span>
              </div>
              <div className="flex gap-4">
                <span className="text-zinc-600 select-none">4</span>
                <span className="text-white">radpresence install</span>
              </div>
              <div className="flex gap-4">
                <span className="text-zinc-600 select-none">5</span>
                <span className="text-white">radpresence start</span>
              </div>
            </div>
          </div>
          <div className="mt-8 flex justify-center">
            <a href={GITHUB_URL + "#getting-started"} target="_blank" rel="noreferrer" className="text-zinc-400 hover:text-white flex items-center gap-2 transition-colors border-b border-zinc-700 hover:border-zinc-400 pb-1">
              <Terminal size={18} />
              Read the full installation guide
            </a>
          </div>
        </section>

      </main>

      {/* Footer */}
      <footer className="border-t border-white/10 bg-zinc-950 py-12 relative z-10">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2 font-display font-medium text-lg">
            <span className="text-blurple-500">RAD</span>
            <span className="text-zinc-400">Presence</span>
          </div>
          <p className="text-zinc-500 text-sm">
            Not officially affiliated with RetroAchievements or Discord.
          </p>
          <div className="flex gap-4">
            <a href={GITHUB_URL} className="text-zinc-500 hover:text-white transition-colors">
              <Github size={20} />
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({ icon, title, description, className = "" }: { icon: React.ReactNode, title: string, description: string, className?: string }) {
  return (
    <div className={`p-8 rounded-3xl border border-white/5 flex flex-col gap-4 group hover:border-white/10 transition-colors ${className}`}>
      <div className="w-12 h-12 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center text-zinc-300 group-hover:text-blurple-400 group-hover:scale-110 transition-all">
        {icon}
      </div>
      <div>
        <h3 className="text-xl font-bold font-display text-white mb-2">{title}</h3>
        <p className="text-zinc-400 leading-relaxed">{description}</p>
      </div>
    </div>
  );
}
