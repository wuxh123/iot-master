import {Component, OnInit, ViewContainerRef} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerService, NzModalService} from "ng-zorro-antd";
import {LinkEditComponent} from "../link-edit/link-edit.component";
import {LinkDetailComponent} from "../link-detail/link-detail.component";

@Component({
  selector: 'app-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss']
})
export class LinkComponent implements OnInit {
  links: [];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;

  statusFilters = [{text: '启动', value: 1}];


  constructor(private as: ApiService, private drawer: NzDrawerService) {
  }

  ngOnInit(): void {
    this.loadFilters();
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('links', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keyword: this.keyword,
    }).subscribe(res => {

      this.links = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }

  loadFilters(): void {
    // this.as.get('distinct/copy/host').subscribe(res => {
    //   console.log('res', res);
    //   this.hosts = res.data.map(h => {
    //     return {
    //       text: h.host,
    //       value: h.host
    //     };
    //   });
    // }, error => {
    //   console.log('error', error);
    // });
  }

  edit(l): void {
    this.drawer.create({
      nzTitle: '编辑',
      nzMaskClosable: false,
      nzWidth: 500,
      nzContent: LinkEditComponent,
      nzContentParams: {
        link: l
      }
    });
  }

  detail(l): void {
    this.drawer.create({
      nzTitle: '详情',
      // nzWidth: 400,
      nzContent: LinkDetailComponent,
      nzContentParams: {
        link: l
      }
    });
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
