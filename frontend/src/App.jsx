
import { useState, useEffect } from 'react';
import './App.css';
import './index.css';

const statusSteps = [
  {
    title: 'Initializing ZK Coprocessor',
    description: 'Starting Brevis core engine and preparing the proof environment.'
  },
  {
    title: 'Establishing Data Provenance',
    description: 'Securely fetching transaction data from Ethereum and verifying data integrity.'
  },
  {
    title: 'Executing Programmable ZK Logic',
    description: 'Running your custom Go circuit to analyze on-chain data using the ZK engine.'
  },
  {
    title: 'Privacy-Preserving Computation',
    description: 'Generating a cryptographic ZK proof that validates computation without revealing raw data.'
  },
  {
    title: 'Finalizing Proof...',
    description: 'Proof is ready to be verified on-chain.'
  }
];

function App() {
  const [txHash, setTxHash] = useState('0x4624b5300e2b26002982adb0b64e643a6056fa0a012a6f4a1f11c5e7b2374e2a');
  const [status, setStatus] = useState('idle');
  const [job, setJob] = useState(null);
  const [currentStep, setCurrentStep] = useState(0);

  useEffect(() => {
    if (status !== 'loading') return;
    const runProcess = async () => {
      const provePromise = fetch('https://silver-fiesta-xjrw6r57gxwcp4pv-3001.app.github.dev/prove', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ txHash }),
      });
      const storyInterval = setInterval(() => {
        setCurrentStep(prevStep => prevStep + 1);
      }, 2500);
      await new Promise(resolve => setTimeout(resolve, statusSteps.length * 2500 + 500));
      clearInterval(storyInterval);
      try {
        const response = await provePromise;
        const data = await response.json();
        if (!response.ok || data.status !== 'success') {
          throw new Error(data.message || 'Proof generation failed');
        }
        setJob(data);
        setStatus('success');
      } catch (error) {
        setJob({ message: error.message });
        setStatus('error');
      }
    };
    runProcess();
  }, [status, txHash]);

  const handleSubmit = (event) => {
    event.preventDefault();
    setCurrentStep(0);
    setJob(null);
    setStatus('loading');
  };
  const isLoading = status === 'loading';

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Brevis ZK Proof Generator</h1>
        <p>Enter a transaction hash to generate a ZK proof.</p>
      </header>
      <main>
        <form className="proof-form" onSubmit={handleSubmit}>
          <input type="text" value={txHash} onChange={(e) => setTxHash(e.target.value)} disabled={isLoading} />
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Generating...' : 'Generate Proof'}
          </button>
        </form>
        <div className={`status-area ${status === 'success' ? 'status-success' : ''} ${status === 'error' ? 'status-error' : ''}`}>
          {status === 'idle' && <p>Please enter a transaction hash and click "Generate Proof".</p>}
          {isLoading && (
            <div>
              <p className="step-title">{statusSteps[Math.min(currentStep, statusSteps.length - 1)].title}</p>
              <p className="step-description">{statusSteps[Math.min(currentStep, statusSteps.length - 1)].description}</p>
              <div className="progress-bar" style={{ marginTop: '20px' }}>
                <div className="progress-bar-inner"></div>
              </div>
            </div>
          )}
          {status === 'success' && (
            <div>
              <h2>Success! ✅</h2>
              <p>{job?.message}</p>
            </div>
          )}
          {status === 'error' && (
            <div>
              <h2>Error ❌</h2>
              <p>{job?.message}</p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;
