'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { apiClient } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { BookOpen, Heart, GraduationCap, Shield, ArrowLeft, Check } from 'lucide-react';

type UserType = 'famille' | 'enseignant' | 'administrateur';

export default function RegisterPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const typeFromUrl = searchParams.get('type') as UserType;
  
  const [selectedType, setSelectedType] = useState<UserType | null>(typeFromUrl || null);
  const [step, setStep] = useState<'type' | 'form'>(typeFromUrl ? 'form' : 'type');
  
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
    phone_number: '',
    // Champs spécifiques selon le type
    family_name: '', // Pour famille
    specialization: '', // Pour enseignant
    qualifications: '', // Pour enseignant
    experience_years: '', // Pour l'affichage seulement
    education_level: '', // Pour l'affichage seulement
    children_ages: '', // Pour l'affichage seulement
    address: '', // Pour l'affichage seulement
    department: '', // Pour l'affichage seulement
  });
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleTypeSelect = (type: UserType) => {
    setSelectedType(type);
    setStep('form');
    // Mettre à jour l'URL sans reload
    const url = new URL(window.location.href);
    url.searchParams.set('type', type);
    window.history.pushState({}, '', url.toString());
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Les mots de passe ne correspondent pas');
      setLoading(false);
      return;
    }

    try {
      // Mapper les types d'utilisateurs
      const roleMapping = {
        famille: 'famille' as const,
        enseignant: 'enseignant' as const,
        administrateur: 'administrator' as const
      };

      const registrationData: {
        username: string;
        email: string;
        password: string;
        role: 'famille' | 'enseignant' | 'administrator';
        phone_number?: string;
        family_name?: string;
        specialization?: string;
        qualifications?: string;
      } = {
        username: formData.username,
        email: formData.email,
        password: formData.password,
        role: roleMapping[selectedType!],
        phone_number: formData.phone_number || undefined,
      };

      // Ajouter les champs spécifiques selon le type
      if (selectedType === 'famille') {
        registrationData.family_name = formData.family_name;
      } else if (selectedType === 'enseignant') {
        registrationData.specialization = formData.specialization;
        // Construire les qualifications à partir des champs du formulaire
        const qualificationsParts = [];
        if (formData.education_level) qualificationsParts.push(`Niveau: ${formData.education_level}`);
        if (formData.experience_years) qualificationsParts.push(`Expérience: ${formData.experience_years}`);
        registrationData.qualifications = qualificationsParts.join(', ');
      }

      const response = await apiClient.register(registrationData);
      
      if (response.error) {
        setError(response.error);
      } else {
        router.push('/dashboard');
        router.refresh();
      }
    } catch (err) {
      setError('Erreur lors de l\'inscription');
    } finally {
      setLoading(false);
    }
  };

  const getTypeConfig = (type: UserType) => {
    const configs = {
      famille: {
        icon: Heart,
        title: 'Famille',
        description: 'Trouvez les meilleurs enseignants pour vos enfants',
        color: 'blue',
        features: [
          'Accès à tous les enseignants qualifiés',
          'Planification flexible des cours',
          'Suivi des progrès en temps réel',
          'Paiements sécurisés'
        ]
      },
      enseignant: {
        icon: GraduationCap,
        title: 'Enseignant',
        description: 'Partagez vos connaissances et développez votre activité',
        color: 'green',
        features: [
          'Créez votre profil professionnel',
          'Gérez vos missions et horaires',
          'Accédez aux ressources pédagogiques',
          'Recevez vos paiements rapidement'
        ]
      },
      administrateur: {
        icon: Shield,
        title: 'Administrateur',
        description: 'Gérez et supervisez la plateforme',
        color: 'purple',
        features: [
          'Dashboard complet de gestion',
          'Supervision des utilisateurs',
          'Validation des rapports',
          'Support et assistance'
        ]
      }
    };
    return configs[type];
  };

  if (step === 'type') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-12">
            <Link href="/" className="inline-flex items-center text-indigo-600 hover:text-indigo-700 mb-8">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Retour à l'accueil
            </Link>
            <div className="flex items-center justify-center space-x-2 mb-6">
              <BookOpen className="h-10 w-10 text-indigo-600" />
              <span className="text-3xl font-bold text-gray-900">EduConnect</span>
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              Choisissez votre profil
            </h1>
            <p className="text-xl text-gray-600">
              Sélectionnez le type de compte qui correspond à vos besoins
            </p>
          </div>

          {/* Type Selection Cards */}
          <div className="grid md:grid-cols-3 gap-8">
            {(['famille', 'enseignant', 'administrateur'] as UserType[]).map((type) => {
              const config = getTypeConfig(type);
              const Icon = config.icon;
              
              return (
                <div 
                  key={type}
                  className="bg-white rounded-2xl p-8 shadow-lg hover:shadow-xl transition-all duration-300 border border-gray-200/50 cursor-pointer group"
                  onClick={() => handleTypeSelect(type)}
                >
                  <div className={`w-16 h-16 bg-gradient-to-br from-${config.color}-500 to-${config.color}-600 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform`}>
                    <Icon className="h-8 w-8 text-white" />
                  </div>
                  
                  <h3 className="text-2xl font-semibold text-gray-900 mb-3 text-center">
                    {config.title}
                  </h3>
                  
                  <p className="text-gray-600 mb-6 text-center">
                    {config.description}
                  </p>
                  
                  <ul className="space-y-3 mb-8">
                    {config.features.map((feature, index) => (
                      <li key={index} className="flex items-center text-gray-600">
                        <Check className="h-4 w-4 text-green-500 mr-3 flex-shrink-0" />
                        {feature}
                      </li>
                    ))}
                  </ul>
                  
                  <Button 
                    className={`w-full bg-${config.color}-600 hover:bg-${config.color}-700`}
                    onClick={() => handleTypeSelect(type)}
                  >
                    Choisir ce profil
                  </Button>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    );
  }

  // Form Step
  const config = selectedType ? getTypeConfig(selectedType) : null;
  const Icon = config?.icon || BookOpen;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <button 
            onClick={() => setStep('type')}
            className="inline-flex items-center text-indigo-600 hover:text-indigo-700 mb-6"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Changer de type de profil
          </button>
          
          <div className={`w-16 h-16 bg-gradient-to-br from-${config?.color}-500 to-${config?.color}-600 rounded-2xl flex items-center justify-center mx-auto mb-4`}>
            <Icon className="h-8 w-8 text-white" />
          </div>
          
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Inscription {config?.title}
          </h1>
          <p className="text-gray-600">
            {config?.description}
          </p>
        </div>

        {/* Form */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}
          
            {/* Informations de base */}
          <div className="space-y-4">
              <h3 className="text-lg font-semibold text-gray-900">Informations personnelles</h3>
              
              <div className="grid md:grid-cols-2 gap-4">
            <div>
                  <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
                    Nom d'utilisateur *
              </label>
              <input
                id="username"
                name="username"
                type="text"
                required
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                value={formData.username}
                onChange={handleChange}
              />
            </div>

            <div>
                  <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                    Adresse email *
              </label>
              <input
                id="email"
                name="email"
                type="email"
                required
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                value={formData.email}
                onChange={handleChange}
              />
                </div>
            </div>

            <div>
                <label htmlFor="phone_number" className="block text-sm font-medium text-gray-700 mb-1">
                  Téléphone
                </label>
                <input
                  id="phone_number"
                  name="phone_number"
                  type="tel"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  value={formData.phone_number}
                  onChange={handleChange}
                />
              </div>

              <div className="grid md:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                    Mot de passe *
              </label>
              <input
                id="password"
                name="password"
                type="password"
                required
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                value={formData.password}
                onChange={handleChange}
              />
            </div>
                
                <div>
                  <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
                    Confirmer le mot de passe *
                  </label>
                  <input
                    id="confirmPassword"
                    name="confirmPassword"
                    type="password"
                    required
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.confirmPassword}
                    onChange={handleChange}
                  />
                </div>
              </div>
            </div>

            {/* Champs spécifiques selon le type */}
            {selectedType === 'enseignant' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-gray-900">Informations professionnelles</h3>
                
                <div>
                  <label htmlFor="specialization" className="block text-sm font-medium text-gray-700 mb-1">
                    Spécialisation
                  </label>
                  <input
                    id="specialization"
                    name="specialization"
                    type="text"
                    placeholder="Ex: Mathématiques, Français, Anglais..."
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.specialization}
                    onChange={handleChange}
                  />
                </div>

                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="experience_years" className="block text-sm font-medium text-gray-700 mb-1">
                      Années d'expérience
                    </label>
                    <select
                      id="experience_years"
                      name="experience_years"
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                      value={formData.experience_years}
                      onChange={handleChange}
                    >
                      <option value="">Sélectionner</option>
                      <option value="0-1">0-1 an</option>
                      <option value="1-3">1-3 ans</option>
                      <option value="3-5">3-5 ans</option>
                      <option value="5-10">5-10 ans</option>
                      <option value="10+">10+ ans</option>
                    </select>
                  </div>
                  
                  <div>
                    <label htmlFor="education_level" className="block text-sm font-medium text-gray-700 mb-1">
                      Niveau d'études
                    </label>
                    <select
                      id="education_level"
                      name="education_level"
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                      value={formData.education_level}
                      onChange={handleChange}
                    >
                      <option value="">Sélectionner</option>
                      <option value="Bac">Bac</option>
                      <option value="Bac+2">Bac+2</option>
                      <option value="Bac+3">Bac+3</option>
                      <option value="Bac+5">Bac+5</option>
                      <option value="Doctorat">Doctorat</option>
                    </select>
                  </div>
                </div>
              </div>
            )}

            {selectedType === 'famille' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-gray-900">Informations familiales</h3>
                
                <div>
                  <label htmlFor="family_name" className="block text-sm font-medium text-gray-700 mb-1">
                    Nom de famille
                  </label>
                  <input
                    id="family_name"
                    name="family_name"
                    type="text"
                    placeholder="Nom de la famille"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.family_name}
                    onChange={handleChange}
                  />
                </div>
                
                <div>
                  <label htmlFor="children_ages" className="block text-sm font-medium text-gray-700 mb-1">
                    Âges des enfants
                  </label>
                  <input
                    id="children_ages"
                    name="children_ages"
                    type="text"
                    placeholder="Ex: 8, 12, 15 ans"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.children_ages}
                    onChange={handleChange}
                  />
          </div>

          <div>
                  <label htmlFor="address" className="block text-sm font-medium text-gray-700 mb-1">
                    Adresse
                  </label>
                  <textarea
                    id="address"
                    name="address"
                    rows={3}
                    placeholder="Votre adresse complète"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.address}
                    onChange={handleChange}
                  />
                </div>
              </div>
            )}

            {selectedType === 'administrateur' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-gray-900">Informations administratives</h3>
                
                <div>
                  <label htmlFor="department" className="block text-sm font-medium text-gray-700 mb-1">
                    Département
                  </label>
                  <select
                    id="department"
                    name="department"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    value={formData.department}
                    onChange={handleChange}
                  >
                    <option value="">Sélectionner</option>
                    <option value="Gestion utilisateurs">Gestion utilisateurs</option>
                    <option value="Support technique">Support technique</option>
                    <option value="Validation académique">Validation académique</option>
                    <option value="Paiements">Paiements</option>
                    <option value="Direction">Direction</option>
                  </select>
                </div>
              </div>
            )}

            <div className="pt-6">
              <Button 
              type="submit"
              disabled={loading}
                className={`w-full bg-${config?.color}-600 hover:bg-${config?.color}-700`}
                size="lg"
            >
                {loading ? 'Inscription en cours...' : 'Créer mon compte'}
              </Button>
          </div>

            <div className="text-center pt-4">
              <Link 
              href="/login"
                className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Déjà un compte ? Se connecter
              </Link>
          </div>
        </form>
        </div>
      </div>
    </div>
  );
} 