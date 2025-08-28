import { TestBed } from '@angular/core/testing';

import { Joke } from './joke';

describe('Joke', () => {
  let service: Joke;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Joke);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
