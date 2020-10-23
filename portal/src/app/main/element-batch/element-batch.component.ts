import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzTableQueryParams} from 'ng-zorro-antd';
import {Router} from '@angular/router';

@Component({
  selector: 'app-element-batch',
  templateUrl: './element-batch.component.html',
  styleUrls: ['./element-batch.component.scss']
})
export class ElementBatchComponent implements OnInit {
  @Input() element: any = {};

  batches: [];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;

  constructor(private as: ApiService, private router: Router) {
  }

  ngOnInit(): void {
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('element/' + this.element.id + '/batches', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keyword: this.keyword,
    }).subscribe(res => {

      this.batches = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }

  create(): void {
    this.router.navigate(['/admin/element/' + this.element.id + '/batch/create']);
  }

  edit(c): void {
    this.router.navigate(['/admin/element/' + this.element.id + '/batch/' + c.id + '/edit']);
  }

  onTableQuery(params: NzTableQueryParams): void {
    const {pageSize, pageIndex, sort, filter} = params;
    this.pageSize = pageSize;
    this.pageIndex = pageIndex;
    const currentSort = sort.find(item => item.value !== null);
    this.sortField = (currentSort && currentSort.key) || null;
    this.sortOrder = (currentSort && currentSort.value) || null;
    this.filters = filter;
    this.load();
  }

  search(): void {
    this.pageIndex = 1;
    this.load();
  }
}
