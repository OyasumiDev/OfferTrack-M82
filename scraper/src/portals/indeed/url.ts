export interface IndeedQuery {
  puesto: string;
  ciudad?: string;
  estado?: string;
  page?: number;
  radiusKm?: number;
  salaryMinAnnualMXN?: number;
  fromAgeDays?: number;
  modality?: string;
  fullTime?: boolean;
  sortByDate?: boolean;
}

const INDEED_BASE_URL = 'https://mx.indeed.com/jobs';

export function BuildIndeedURL(q: IndeedQuery): string {
  const params = new URLSearchParams();
  params.set('q', q.puesto);

  if (q.modality === 'remote') {
    params.set('l', 'remote');
  } else if (q.ciudad && q.estado) {
    params.set('l', `${q.ciudad}, ${q.estado}`);
  } else if (q.ciudad) {
    params.set('l', q.ciudad);
  } else if (q.estado) {
    params.set('l', q.estado);
  }

  if (q.radiusKm && q.radiusKm !== 25) {
    params.set('radius', String(q.radiusKm));
  }

  if (q.salaryMinAnnualMXN && q.salaryMinAnnualMXN > 0) {
    params.set('salaryType', '$' + q.salaryMinAnnualMXN.toLocaleString('en-US'));
  }

  if (q.fromAgeDays && q.fromAgeDays > 0) {
    params.set('fromage', String(q.fromAgeDays));
  }

  if (q.page && q.page > 1) {
    params.set('start', String((q.page - 1) * 10));
  }

  if (q.sortByDate) {
    params.set('sort', 'date');
  }

  const attrs: string[] = [];
  if (q.modality === 'hybrid') attrs.push('PAXZC');
  if (q.fullTime) attrs.push('CF3CP');
  if (attrs.length > 0) {
    params.set('sc', '0kf:' + attrs.map(a => `attr(${a})`).join('') + ';');
  }

  return INDEED_BASE_URL + '?' + params.toString();
}
