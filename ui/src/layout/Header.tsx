import { InfoIcon } from 'lucide-react';
import { Link } from 'react-router-dom';


export default function Header() {
  return (
    <div className="border-b border-gray-200">
      <div className="mx-auto flex max-w-7xl items-center justify-between py-4">
        <Link to="/">
          <div className="flex items-center gap-4">
            <img src="../assets/mentis-logo.png" className="h-10 w-10" ></img>
            <span className="text-2xl font-bold text-emerald-600">Mentis</span>
          </div>
        </Link>
        <Link to="/help">
          <div className="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700">
            How to Record
            <InfoIcon className="h-5 w-5" />
          </div>
        </Link>
      </div>
    </div>
  );
}
